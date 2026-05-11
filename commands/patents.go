package commands

import "github.com/spf13/cobra"

var patentsCmd = &cobra.Command{
	Use:   "patents",
	Short: "",
}

func init() {
	rootCmd.AddCommand(patentsCmd)
}
