package main

import "github.com/kraytonsoft/daxa/sdk/daxa"

func Handler(ctx daxa.RequestContext) (daxa.Response, error) {
	return daxa.Response{
		Status: 200,
		Body:   []byte("Hello from Daxa"),
	}, nil
}
