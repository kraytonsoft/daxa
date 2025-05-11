package daxa

import "net/http"

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
