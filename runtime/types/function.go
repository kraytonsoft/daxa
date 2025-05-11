package types

import "net/http"

type Function struct {
	ID         string `json:"id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	PluginPath string `json:"plugin"` // runtime-generated
}

type Manifest struct {
	Functions []Function `json:"functions"`
}

type RequestContext struct {
	Request *http.Request
	Writer  http.ResponseWriter
}

type Response struct {
	Status  int
	Headers map[string]string
	Body    []byte
}

type DaxaFunc func(ctx RequestContext) (Response, error)
