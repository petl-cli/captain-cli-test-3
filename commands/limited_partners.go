package commands

import "github.com/spf13/cobra"

var limitedPartnersCmd = &cobra.Command{
	Use:   "limited-partners",
	Short: "",
}

func init() {
	rootCmd.AddCommand(limitedPartnersCmd)
}
