package structs

type UrlConfig struct {
	Method string
	Path   string // (its a regex)
}

type Request struct {
	Params map[string]interface{}
	Body   map[string]interface{}
	Query  map[string]interface{}
}
