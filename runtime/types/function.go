package types

type Function struct {
	ID         string `json:"id"`
	Method     string `json:"method"`
	Path       string `json:"path"`
	PluginPath string `json:"plugin"` // runtime-generated
}

type Manifest struct {
	Functions []Function `json:"functions"`
}
