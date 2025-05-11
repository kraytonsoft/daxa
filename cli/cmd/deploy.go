// 📁 cli/cmd/deploy.go
package cmd

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path]",
	Short: "Deploy a Daxa function by zipping the source folder and sending it to the runtime",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		srcPath := args[0]

		// Validate required files
		required := []string{"main.go", "go.mod", "daxa.json"}
		for _, f := range required {
			if _, err := os.Stat(filepath.Join(srcPath, f)); errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("missing required file: %s", f)
			}
		}

		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)

		err := filepath.Walk(srcPath, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			rel, _ := filepath.Rel(srcPath, p)
			f, _ := zw.Create(rel)
			src, _ := os.Open(p)
			defer src.Close()
			_, err = io.Copy(f, src)
			return err
		})
		if err != nil {
			return fmt.Errorf("failed to zip: %w", err)
		}
		zw.Close()

		resp, err := http.Post("http://localhost:36365/deploy/source", "application/zip", bytes.NewReader(zipBuf.Bytes()))
		if err != nil {
			return fmt.Errorf("failed to deploy: %w", err)
		}
		defer resp.Body.Close()
		io.Copy(os.Stdout, resp.Body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
