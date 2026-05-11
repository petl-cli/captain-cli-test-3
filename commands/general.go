package commands

import "github.com/spf13/cobra"

var generalCmd = &cobra.Command{
	Use:   "general",
	Short: "",
}

func init() {
	rootCmd.AddCommand(generalCmd)
}
