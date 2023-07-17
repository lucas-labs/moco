package http

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	DELETE  HttpMethod = "DELETE"
	HEAD    HttpMethod = "HEAD"
	OPTIONS HttpMethod = "OPTIONS"
	PATCH   HttpMethod = "PATCH"
)

func (h HttpMethod) String() string {
	return string(h)
}

func (h HttpMethod) IsValid() bool {
	switch h {
	case GET, POST, PUT, DELETE, HEAD, OPTIONS, PATCH:
		return true
	}
	return false
}

func IsValidHttpStatus(status int) bool {
	return status >= 100 && status < 600
}
