package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "daxa",
	Short: "Daxagrid CLI",
	Long:  "CLI for interacting with the Daxagrid runtime cluster.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
