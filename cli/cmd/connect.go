package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [name] [host]",
	Short: "Connect to a Daxa runtime and save the kontext",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		host := args[1]
		err := SaveContext(name, host)
		if err != nil {
			return err
		}
		fmt.Printf("[OK] - Connected to %s (%s)\n", name, host)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}
