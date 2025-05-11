package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	kontext "daxa/cli/kontext"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:   "deploy [path]",
	Short: "Deploy a Go function project to the Daxa runtime",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		zipBuf := new(bytes.Buffer)
		zw := zip.NewWriter(zipBuf)

		err := filepath.Walk(args[0], func(p string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}

			rel, _ := filepath.Rel(args[0], p)
			f, _ := zw.Create(rel)
			src, err := os.Open(p)
			if err != nil {
				return err
			}
			defer src.Close()
			_, err = io.Copy(f, src)
			return err
		})
		if err != nil {
			return err
		}

		zw.Close()

		host, err := kontext.GetCurrentHost()
		if err != nil {
			return fmt.Errorf("no connected runtime: run `daxa connect` first")
		}

		url := fmt.Sprintf("%s/deploy/source", host)
		resp, err := http.Post(url, "application/zip", bytes.NewReader(zipBuf.Bytes()))

		io.Copy(os.Stdout, resp.Body)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
