package commands

import "github.com/spf13/cobra"

var datasetsCmd = &cobra.Command{
	Use:   "datasets",
	Short: "",
}

func init() {
	rootCmd.AddCommand(datasetsCmd)
}
