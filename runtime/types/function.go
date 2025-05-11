package types

type Manifest struct {
	Functions []Function `json:"functions"`
}

// Function defines a deployable Go plugin handler
type Function struct {
	ID         string `json:"id"`     // e.g. "hello"
	Method     string `json:"method"` // e.g. "GET"
	Path       string `json:"path"`   // e.g. "/hello"
	PluginPath string `json:"plugin"` // e.g. "hello.so"
}
