package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strconv"

	"moco/config"
	"moco/server/cors"
	"moco/utils/compare"
	"moco/utils/logger"
)

//go:embed static/favicon.svg
var favicon []byte

var log = logger.New("server")

func Start(
	port int,
	allawedOrigins []string,
	endpoints []config.Endpoint,
) error {
	fatalCh := make(chan error)

	url := "127.0.0.1:" + strconv.Itoa(port)
	mux := http.NewServeMux()
	mux.HandleFunc("/", makeHandler(endpoints))

	handler := cors.New(cors.Options{
		AllowedOrigins: allawedOrigins,
	}).Handler(mux)

	go serve(handler, url, fatalCh)
	return <-fatalCh
}

func serve(handler http.Handler, url string, fatalCh chan error) {
	listener, err := net.Listen("tcp", url)
	if err != nil {
		fatalCh <- fmt.Errorf("failed to listen on %s: %w", url, err)
		return
	}

	log.Infof("Listening on http://%s", url)

	err = http.Serve(listener, handler)
	if err != nil {
		fatalCh <- err
	}
}

func makeHandler(endpoints []config.Endpoint) func(http.ResponseWriter, *http.Request) {
	for _, endpoint := range endpoints {
		log.Infof("Registered Endpoint %s %s", fmt.Sprintf("%-5s", endpoint.Url.Method.String()), endpoint.Url.Path)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/favicon.ico" {
			faviconHandler(w)
			return
		}

		log.Info(r.Method, r.URL.Path)

		for _, endpoint := range endpoints {
			if endpoint.Url.Matches(r.Method, r.URL.Path) {
				req, err := GetRequest(r, endpoint.Url)

				if err != nil {
					log.Errore(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				endpointHandler(w, r, req, endpoint.Responses)
				return
			}
		}

		http.NotFound(w, r)
	}
}

func endpointHandler(w http.ResponseWriter, r *http.Request, request Request, responses []config.ResponseConfig) {
	// first we check if the request.Params, request.Body and request.Query matches any of the responses[n].Request.Params, responses[n].Request.Body and responses[n].Request.Query
	// all of them must match, and they do match, we return the response specified in the responses[n]
	// if none of them match, we check if any of the responses[n].Request has IsDefault set to true, and if so, we return that response
	// if none of them match, we return a 404
	var defaultResponse *config.ResponseConfig

	for _, response := range responses {
		if response.Request.IsDefault {
			// we save the default response in case we need it later,
			// so we don't have to loop again
			defaultResponse = &response
			continue
		}

		log.Info("Checking response")

		// check if request.Params matches response.Request.Params deeply
		if compare.DeepEqual(response.Request.Params, request.Params) &&
			compare.DeepEqual(response.Request.Body, request.Body) &&
			compare.DeepEqual(response.Request.Query, request.Query) {
			sendResponse(w, response)
			return
		} else {
			// Perform actions if any of the fields don't match
			continue
		}
	}

	if defaultResponse != nil {
		sendResponse(w, *defaultResponse)
		return
	}

	// return 404
	http.NotFound(w, r)
}

func sendResponse(w http.ResponseWriter, response config.ResponseConfig) {
	headersToResponse := response.ResponseHeaders
	statusCodeToResponse := response.Status

	var bodyToResponse []byte
	var marshableResponseBody = make(map[string]interface{})
	for key, value := range response.ResponseBody {
		marshableResponseBody[key.(string)] = value
	}
	bodyToResponse, err := json.Marshal(marshableResponseBody)

	if err != nil {
		log.Errore(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// set headers
	for key, value := range headersToResponse {
		w.Header().Set(key, value)
	}

	// check that we have set a Content-Type header
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	// send response
	w.WriteHeader(statusCodeToResponse)
	w.Write(bodyToResponse)
}

func faviconHandler(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "image/svg+xml")
	w.Write(favicon)
}

type UrlConfig struct {
	Method string
	Path   string // (its a regex)
}

type Request struct {
	Params map[string]interface{}
	Body   map[string]interface{}
	Query  map[string]interface{}
}

func GetRequest(r *http.Request, urlCfg config.UrlConfig) (Request, error) {
	var request Request

	// Parsing URL parameters
	regex := regexp.MustCompile(urlCfg.Path)
	match := regex.FindStringSubmatch(r.URL.Path)

	if len(match) > 1 {
		params := make(map[string]interface{})
		for i, name := range regex.SubexpNames() {
			if i != 0 && name != "" {
				params[name] = match[i]
			}
		}
		request.Params = params
	}

	// Parsing request body
	if r.ContentLength > 0 {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return request, err
		}
		err = json.Unmarshal(body, &request.Body)
		if err != nil {
			return request, err
		}
	}

	// Parsing URL query parameters
	queryParams := r.URL.Query()
	if len(queryParams) > 0 {
		// assign queryParams to request.Query
		request.Query = make(map[string]interface{})
		for key, value := range queryParams {
			request.Query[key] = value
		}
	}

	return request, nil
}
