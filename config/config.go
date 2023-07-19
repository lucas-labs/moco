package config

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	httpUtils "moco/utils/http"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name        string     `yaml:"name"`
	Version     string     `yaml:"version"`
	Port        int        `yaml:"port"`
	Cors        []string   `yaml:"cors"`
	Description string     `yaml:"description"`
	Endpoints   []Endpoint `yaml:"endpoints"`
}

type Endpoint struct {
	Url         UrlConfig
	SideEffects []SideEffectConfig `yaml:"sideEffects"`
	Responses   []ResponseConfig
	Childs      []Endpoint
}

type SideEffectConfig struct {
	Type   string                 `yaml:"type"`
	Entity string                 `yaml:"entity"`
	Data   map[string]interface{} `yaml:"data"`
}

type ResponseConfig struct {
	Request         Request
	Status          int
	ResponseBody    map[string]interface{}
	ResponseHeaders map[string]string
}

func (r *ResponseConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("invalid response config, expected array, example:  [{ params: { id: \"2\" } }, 200, { id: 2, name: \"Jane Doe\" }]")
	}

	if len(value.Content) != 3 && len(value.Content) != 4 {
		return fmt.Errorf("invalid response config, expected array with 3 or 4 elements, example: [{ params: { id: \"2\" } }, 200, { id: 2, name: \"Jane Doe\" }]")
	}

	rawReq := value.Content[0]
	rawStatus := value.Content[1]
	rawBody := value.Content[2]
	var rawHeaders *yaml.Node

	if len(value.Content) == 4 {
		rawHeaders = value.Content[3]
	} else {
		rawHeaders = nil
	}

	var req Request
	var status int
	var body map[string]interface{}
	var headers map[string]string

	if rawReq.Kind != yaml.MappingNode && rawReq.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid request at response config, expected object or \"*\", example: { params: { id: \"2\" } })")
	}

	if rawStatus.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid status at response config, expected number, example: 200")
	}

	if rawBody.Kind != yaml.MappingNode && rawBody.Kind != yaml.SequenceNode {
		return fmt.Errorf("invalid body at response config, expected object or array, example: { id: 2, name: \"Jane Doe\" }")
	}

	if rawHeaders != nil && rawHeaders.Kind != yaml.MappingNode {
		return fmt.Errorf("invalid headers at response config, expected object, example: { \"Content-Type\": \"application/json\" }")
	}

	// decode final values

	if rawReq.Kind == yaml.MappingNode {
		err := rawReq.Decode(&req)
		if err != nil {
			return err
		}
	}

	if rawReq.Kind == yaml.ScalarNode {
		// check it's "*"
		if rawReq.Value != "*" {
			return fmt.Errorf("invalid request at response config, expected object or \"*\", example: { params: { id: \"2\" } })")
		}

		// if it's "*", set IsDefault to true
		req = Request{IsDefault: true}
	}

	// convert rawStatus.Value from string to int
	status, err := strconv.Atoi(rawStatus.Value)
	if err != nil {
		return fmt.Errorf("invalid status at response config, expected number, example: 200")
	} else if !httpUtils.IsValidHttpStatus(status) {
		return fmt.Errorf("invalid status at response config, expected valid http status, example: 200")
	}

	err = rawBody.Decode(&body)
	if err != nil {
		return err
	}

	if rawHeaders != nil {
		err = rawHeaders.Decode(&headers)
		if err != nil {
			return err
		}
	}

	r.Request = req
	r.Status = status
	r.ResponseBody = body
	r.ResponseHeaders = headers

	return nil
}

type Request struct {
	IsDefault  bool
	Params     map[string]interface{} `yaml:"params"`
	Body       map[string]interface{} `yaml:"body"`
	Query      map[string]interface{} `yaml:"query"`
	IsLoggedIn bool                   `yaml:"isLoggedIn"`
}

type UrlConfig struct {
	Method httpUtils.HttpMethod `yaml:"method"`
	Path   string               `yaml:"path"`
}

func (u *UrlConfig) Matches(method string, path string) bool {
	// remove prefix and trailing slash if they exist
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	uPath := u.Path
	if len(uPath) > 1 && uPath[len(uPath)-1] == '/' {
		uPath = uPath[:len(uPath)-1]
	}
	if len(uPath) > 0 && uPath[0] == '/' {
		uPath = uPath[1:]
	}

	// if uPath regex doesn't start with ^, and doesn't end with $, add them
	if len(uPath) > 0 && uPath[0] != '^' {
		uPath = "^" + uPath
	}
	if len(uPath) > 0 && uPath[len(uPath)-1] != '$' {
		uPath = uPath + "$"
	}

	// path is a regex
	return u.Method == httpUtils.HttpMethod(method) && regexp.MustCompile(uPath).MatchString(path)
}

func (o *UrlConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.SequenceNode {
		return fmt.Errorf("invalid url config, expected array, example: [GET, /path]")
	}

	if len(value.Content) != 2 {
		return fmt.Errorf("invalid url config, expected array with 2 elements, example: [GET, /path]")
	}

	method := value.Content[0]
	path := value.Content[1]

	if method.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid url config, expected method to be a string, example: [GET, /path]")
	}

	if path.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid url config, expected path to be a string, example: [GET, /path]")
	}

	o.Method = httpUtils.HttpMethod(method.Value)

	if !o.Method.IsValid() {
		return fmt.Errorf("invalid url config, expected method to be a valid http method, example: GET or POST")
	}

	o.Path = path.Value

	return nil
}

func Load() (*Config, string, error) {
	var config Config
	configPaths := []string{
		"config.yml",
		"config.yaml",
		"data.yml",
		"data.yaml",
	}

	// check if one of the config files exists
	var configPath string
	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			break
		}
	}

	if configPath == "" {
		return nil, "", fmt.Errorf("no config file found, expected one of: %v", configPaths)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, configPath, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, configPath, err
	}

	return &config, configPath, nil
}
