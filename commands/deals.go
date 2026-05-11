package commands

import "github.com/spf13/cobra"

var dealsCmd = &cobra.Command{
	Use:   "deals",
	Short: "",
}

func init() {
	rootCmd.AddCommand(dealsCmd)
}
