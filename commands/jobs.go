package commands

import "github.com/spf13/cobra"

var jobsCmd = &cobra.Command{
	Use:   "jobs",
	Short: "",
}

func init() {
	rootCmd.AddCommand(jobsCmd)
}
