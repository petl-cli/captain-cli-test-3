package commands

import "github.com/spf13/cobra"

var collectionsCmd = &cobra.Command{
	Use:   "collections",
	Short: "",
}

func init() {
	rootCmd.AddCommand(collectionsCmd)
}
