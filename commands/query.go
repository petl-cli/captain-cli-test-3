package commands

import "github.com/spf13/cobra"

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "",
}

func init() {
	rootCmd.AddCommand(queryCmd)
}
