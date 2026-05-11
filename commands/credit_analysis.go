package commands

import "github.com/spf13/cobra"

var creditAnalysisCmd = &cobra.Command{
	Use:   "credit-analysis",
	Short: "",
}

func init() {
	rootCmd.AddCommand(creditAnalysisCmd)
}
