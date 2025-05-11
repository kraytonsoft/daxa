package handlers

import (
	"archive/zip"
	"daxa/runtime/compiler"
	"daxa/runtime/registry"
	"daxa/runtime/types"
	"io"
	_ "io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var fnRegistry *registry.FunctionRegistry

func Init(r *registry.FunctionRegistry) {
	fnRegistry = r
}

func handleDeploySource(w http.ResponseWriter, r *http.Request) {
	tmpDir := filepath.Join(os.TempDir(), "daxa", uuid.New().String())
	os.MkdirAll(tmpDir, 0755)

	zipFile := filepath.Join(tmpDir, "source.zip")
	f, _ := os.Create(zipFile)
	io.Copy(f, r.Body)
	f.Close()

	unzip(zipFile, tmpDir)

	pluginPath, err := compiler.BuildPlugin(tmpDir, "function.so")
	if err != nil {
		http.Error(w, "Failed to compile: "+err.Error(), 500)
		return
	}

	fn := types.Function{
		ID:         "auto",
		Method:     "GET",
		Path:       "/hello", // read this from config later
		PluginPath: pluginPath,
	}

	fnRegistry.Register(fn)

	w.Write([]byte("Deployed + Compiled"))
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
