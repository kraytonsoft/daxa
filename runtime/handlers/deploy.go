package handlers

import (
	"archive/zip"
	"daxa/runtime/compiler"
	"daxa/runtime/registry"
	"daxa/runtime/types"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

var reg *registry.FunctionRegistry

func Init(r *registry.FunctionRegistry) {
	reg = r
}

func HandleDeploySource(w http.ResponseWriter, r *http.Request) {
	buildID := uuid.New().String()
	tmpDir := filepath.Join(os.TempDir(), "daxa", buildID)
	os.MkdirAll(tmpDir, 0755)

	zipFile := filepath.Join(tmpDir, "src.zip")
	f, _ := os.Create(zipFile)
	io.Copy(f, r.Body)
	f.Close()

	unzip(zipFile, tmpDir)

	manifestPath := filepath.Join(tmpDir, "daxa.json")
	manifestFile, err := os.ReadFile(manifestPath)
	if err != nil {
		http.Error(w, "Missing daxa.json", 400)
		return
	}

	var manifest types.Manifest
	if err := json.Unmarshal(manifestFile, &manifest); err != nil {
		http.Error(w, "Invalid manifest", 400)
		return
	}

	// Compile plugin
	pluginPath, err := compiler.BuildPlugin(tmpDir, "function.so")
	if err != nil {
		http.Error(w, "Build error: "+err.Error(), 500)
		return
	}

	// Register each function in manifest
	for _, fn := range manifest.Functions {
		fn.PluginPath = pluginPath
		if err := reg.Register(fn); err != nil {
			http.Error(w, "Register error: "+err.Error(), 500)
			return
		}
		http.HandleFunc(fn.Path, func(w http.ResponseWriter, r *http.Request) {
			if h, ok := reg.Handler(fn.Path); ok {
				h(w, r)
			} else {
				http.Error(w, "Handler not found", 500)
			}
		})
	}

	w.Write([]byte("OK"))
}

func unzip(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0755)
			continue
		}
		rc, _ := f.Open()
		dst, _ := os.Create(path)
		io.Copy(dst, rc)
		dst.Close()
		rc.Close()
	}
	return nil
}
