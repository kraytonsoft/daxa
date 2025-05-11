package compiler

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildPlugin(sourceDir string, outName string) (string, error) {
	outPath := filepath.Join(sourceDir, outName)
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outPath, filepath.Join(sourceDir, "main.go"))
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	cmd.Dir = sourceDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("build failed: %v\nOutput: %s", err, string(output))
	}

	return outPath, nil
}
