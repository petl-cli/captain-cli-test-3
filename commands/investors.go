package commands

import "github.com/spf13/cobra"

var investorsCmd = &cobra.Command{
	Use:   "investors",
	Short: "",
}

func init() {
	rootCmd.AddCommand(investorsCmd)
}
