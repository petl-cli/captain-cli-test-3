package commands

import "github.com/spf13/cobra"

var peopleCmd = &cobra.Command{
	Use:   "people",
	Short: "",
}

func init() {
	rootCmd.AddCommand(peopleCmd)
}
