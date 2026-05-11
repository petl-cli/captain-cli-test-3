package commands

import "github.com/spf13/cobra"

var companiesCmd = &cobra.Command{
	Use:   "companies",
	Short: "",
}

func init() {
	rootCmd.AddCommand(companiesCmd)
}
