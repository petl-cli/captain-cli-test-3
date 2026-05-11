package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rishimantri795/CLICreator/runtime/feedback"
	"github.com/rishimantri795/CLICreator/runtime/telemetry"
	"github.com/spf13/cobra"
)

// feedbackEndpoint and feedbackToken are baked in at generation time.
// An empty token disables the command (returns a clear error).
const (
	feedbackEndpoint = "http://localhost:3000/api/feedback/ingest"
	feedbackToken    = "76eed3d4-cff4-4885-860d-9d0f81c18f49"
)

var feedbackFlags struct {
	about string
}

var feedbackCmd = &cobra.Command{
	Use:   "feedback [message]",
	Short: "Send feedback to the API provider about this CLI",
	Long: `Submit feedback to the API provider who publishes this CLI.

Useful when something is unclear, broken, or missing — agents and humans
can both use this channel. Feedback is sent to the provider's dashboard
along with the CLI version and (if detected) the agent runtime in use.

Examples:
  captain-api-v2 feedback "the --filter regex syntax is unclear"
  captain-api-v2 feedback --about "create-widget" "exit code 0 but no output"
`,
	Args: cobra.MinimumNArgs(1),
	RunE: runFeedback,
}

func init() {
	feedbackCmd.Flags().StringVar(&feedbackFlags.about, "about", "",
		"Optional command this feedback is about, e.g. \"create-widget\"")
	rootCmd.AddCommand(feedbackCmd)
}

func runFeedback(cmd *cobra.Command, args []string) error {
	if feedbackToken == "" {
		return fmt.Errorf("feedback is not enabled for this CLI")
	}

	message := strings.TrimSpace(strings.Join(args, " "))
	if message == "" {
		return fmt.Errorf("message is required")
	}

	caller := telemetry.DetectCaller()

	id, err := feedback.Submit(context.Background(), feedbackEndpoint, feedbackToken, feedback.Payload{
		CLIVersion:     "0.1.4",
		Message:        message,
		CommandContext: feedbackFlags.about,
		AgentType:      caller.AgentType,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "feedback not delivered: %v\n", err)
		return err
	}

	if id != "" {
		fmt.Fprintf(os.Stdout, "feedback submitted (id: %s)\n", id)
	} else {
		fmt.Fprintln(os.Stdout, "feedback submitted")
	}
	return nil
}
