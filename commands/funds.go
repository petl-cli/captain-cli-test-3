package commands

import "github.com/spf13/cobra"

var fundsCmd = &cobra.Command{
	Use:   "funds",
	Short: "",
}

func init() {
	rootCmd.AddCommand(fundsCmd)
}
