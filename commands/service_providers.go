package commands

import "github.com/spf13/cobra"

var serviceProvidersCmd = &cobra.Command{
	Use:   "service-providers",
	Short: "",
}

func init() {
	rootCmd.AddCommand(serviceProvidersCmd)
}
