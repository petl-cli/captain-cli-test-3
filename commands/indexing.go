package commands

import "github.com/spf13/cobra"

var indexingCmd = &cobra.Command{
	Use:   "indexing",
	Short: "",
}

func init() {
	rootCmd.AddCommand(indexingCmd)
}
