package registry

import (
	"daxa/runtime/types"
	"daxa/sdk/daxa"
	"fmt"
	"net/http"
	"plugin"
	"sync"
)

type VersionedHandler struct {
	Version int
	Handler http.HandlerFunc
}

type FunctionRegistry struct {
	mu       sync.RWMutex
	versions map[string][]VersionedHandler
	latest   map[string]http.HandlerFunc
}

func New() *FunctionRegistry {
	return &FunctionRegistry{
		versions: make(map[string][]VersionedHandler),
		latest:   make(map[string]http.HandlerFunc),
	}
}

func (fr *FunctionRegistry) Register(fn types.Function) error {
	p, err := plugin.Open(fn.PluginPath)
	if err != nil {
		return fmt.Errorf("failed to load plugin: %w", err)
	}

	sym, err := p.Lookup("Handler")
	if err != nil {
		return fmt.Errorf("missing Handler symbol: %w", err)
	}

	daxaFn, ok := sym.(func(daxa.RequestContext) (daxa.Response, error))
	if !ok {
		return fmt.Errorf("invalid handler signature")
	}

	httpHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := daxa.RequestContext{
			Request: r,
			Writer:  w,
		}
		resp, err := daxaFn(ctx)
		if err != nil {
			http.Error(w, "Function error", 500)
			return
		}
		for k, v := range resp.Headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(resp.Status)
		w.Write(resp.Body)
	}

	fr.mu.Lock()
	defer fr.mu.Unlock()

	version := len(fr.versions[fn.Path]) + 1
	fr.versions[fn.Path] = append(fr.versions[fn.Path], VersionedHandler{
		Version: version,
		Handler: httpHandler,
	})
	fr.latest[fn.Path] = httpHandler
	return nil
}

func (fr *FunctionRegistry) Handler(path string) (http.HandlerFunc, bool) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()
	h, ok := fr.latest[path]
	return h, ok
}
