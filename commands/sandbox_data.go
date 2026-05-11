package commands

import "github.com/spf13/cobra"

var sandboxDataCmd = &cobra.Command{
	Use:   "sandbox-data",
	Short: "",
}

func init() {
	rootCmd.AddCommand(sandboxDataCmd)
}
