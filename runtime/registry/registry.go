package registry

import (
	"daxa/runtime/types"
	"fmt"
	"net/http"
	"plugin"
	"sync"
)

// VersionedHandler keeps track of a specific version of a function
type VersionedHandler struct {
	Version int
	Handler http.HandlerFunc
}

// FunctionRegistry stores versions and active handler per route
type FunctionRegistry struct {
	mu       sync.RWMutex
	versions map[string][]VersionedHandler
	latest   map[string]http.HandlerFunc
}

// New creates a new empty registry
func New() *FunctionRegistry {
	return &FunctionRegistry{
		versions: make(map[string][]VersionedHandler),
		latest:   make(map[string]http.HandlerFunc),
	}
}

// Register loads a Go plugin and stores it in the registry
func (fr *FunctionRegistry) Register(fn types.Function) error {
	p, err := plugin.Open(fn.PluginPath)
	if err != nil {
		return fmt.Errorf("plugin load error: %w", err)
	}

	sym, err := p.Lookup("Handler")
	if err != nil {
		return fmt.Errorf("plugin symbol 'Handler' not found")
	}

	handler, ok := sym.(func(http.ResponseWriter, *http.Request))
	if !ok {
		return fmt.Errorf("Handler has invalid signature")
	}

	fr.mu.Lock()
	defer fr.mu.Unlock()

	version := len(fr.versions[fn.Path]) + 1

	fr.versions[fn.Path] = append(fr.versions[fn.Path], VersionedHandler{
		Version: version,
		Handler: handler,
	})

	fr.latest[fn.Path] = handler
	return nil
}

// Handler returns the most recent handler for a route
func (fr *FunctionRegistry) Handler(path string) (http.HandlerFunc, bool) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()
	h, ok := fr.latest[path]
	return h, ok
}

// Rollback swaps the active handler to a previous version
func (fr *FunctionRegistry) Rollback(path string, version int) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	versions, ok := fr.versions[path]
	if !ok || version < 1 || version > len(versions) {
		return fmt.Errorf("invalid rollback target")
	}

	fr.latest[path] = versions[version-1].Handler
	return nil
}

// Versions returns the version count for a path
func (fr *FunctionRegistry) Versions(path string) int {
	fr.mu.RLock()
	defer fr.mu.RUnlock()
	return len(fr.versions[path])
}
