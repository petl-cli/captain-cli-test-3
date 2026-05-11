package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/rishimantri795/CLICreator/runtime/httpclient"
	"github.com/rishimantri795/CLICreator/runtime/output"
	"github.com/spf13/cobra"
)

var datasetsSearchMedicalPapersCmd = &cobra.Command{
	Use:   "search-medical-papers",
	Short: "Search Medical Papers",
	RunE:  runDatasetsSearchMedicalPapers,
}

var datasetsSearchMedicalPapersFlags struct {
	question      string
	maxSources    int
	includeTrials bool
	recencyYears  string
	stream        bool
	body          string
}

func init() {
	datasetsSearchMedicalPapersCmd.Flags().StringVar(&datasetsSearchMedicalPapersFlags.question, "question", "", "Natural-language question.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsSearchMedicalPapersCmd.Flags().IntVar(&datasetsSearchMedicalPapersFlags.maxSources, "max-sources", 0, "Target number of cited sources in the final answer.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsSearchMedicalPapersCmd.Flags().BoolVar(&datasetsSearchMedicalPapersFlags.includeTrials, "include-trials", false, "Whether the agent may call ClinicalTrials.gov.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsSearchMedicalPapersCmd.Flags().StringVar(&datasetsSearchMedicalPapersFlags.recencyYears, "recency-years", "", "Prefer evidence within the last N years where the question allows.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsSearchMedicalPapersCmd.Flags().BoolVar(&datasetsSearchMedicalPapersFlags.stream, "stream", false, "If true, response is text/event-stream; otherwise JSON.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsSearchMedicalPapersCmd.Flags().StringVar(&datasetsSearchMedicalPapersFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	datasetsCmd.AddCommand(datasetsSearchMedicalPapersCmd)
}

func runDatasetsSearchMedicalPapers(cmd *cobra.Command, args []string) error {
	// --schema: print full input/output type contract without making any network call.
	if rootFlags.schema {
		type flagSchema struct {
			Name        string `json:"name"`
			Type        string `json:"type"`
			Required    bool   `json:"required"`
			Location    string `json:"location"`
			Description string `json:"description,omitempty"`
		}
		var flags []flagSchema
		flags = append(flags, flagSchema{
			Name:        "question",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Natural-language question.",
		})
		flags = append(flags, flagSchema{
			Name:        "max-sources",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Target number of cited sources in the final answer.",
		})
		flags = append(flags, flagSchema{
			Name:        "include-trials",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Whether the agent may call ClinicalTrials.gov.",
		})
		flags = append(flags, flagSchema{
			Name:        "recency-years",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Prefer evidence within the last N years where the question allows.",
		})
		flags = append(flags, flagSchema{
			Name:        "stream",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "If true, response is text/event-stream; otherwise JSON.",
		})

		type responseSchema struct {
			Status      string `json:"status"`
			ContentType string `json:"content_type,omitempty"`
			Description string `json:"description,omitempty"`
		}
		var responses []responseSchema
		responses = append(responses, responseSchema{
			Status:      "200",
			ContentType: "application/json",
			Description: "Synthesized answer with cited sources. Returns JSON by default; when `stream=true` the response is `text/event-stream` with `tool_use`, `tool_result_summary`, `text_delta`, and `done` events.",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "",
			Description: "Malformed body or input validation error.",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Missing or invalid API key.",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "",
			Description: "API key does not belong to the organization, datasets access disabled, or trial/billing limits exceeded.",
		})
		responses = append(responses, responseSchema{
			Status:      "422",
			ContentType: "application/json",
			Description: "Validation Error",
		})

		schema := map[string]any{
			"command":     "search-medical-papers",
			"description": "Search Medical Papers",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/datasets/scientific/medical/ask",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": true,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"creates_resource"},
				"impact":       "medium",
			},
			"requires_auth": false,
		}
		data, _ := json.MarshalIndent(schema, "", "  ")
		fmt.Fprintln(_stdoutCounter, string(data))
		return nil
	}

	cfg, err := rootConfig()
	if err != nil {
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	client := httpclient.New(cfg.BaseURL, cfg.AuthProvider())
	client.Debug = rootFlags.debug
	client.DryRun = rootFlags.dryRun
	if rootFlags.noRetries {
		client.RetryConfig.MaxRetries = 0
	}

	// Build path params
	pathParams := map[string]string{}

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/datasets/scientific/medical/ask", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if datasetsSearchMedicalPapersFlags.body != "" {
		if err := json.Unmarshal([]byte(datasetsSearchMedicalPapersFlags.body), &bodyMap); err != nil {
			_invState.errorType = "parse_error"
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("question") {
		bodyMap["question"] = datasetsSearchMedicalPapersFlags.question
	}
	if cmd.Flags().Changed("max-sources") {
		bodyMap["max_sources"] = datasetsSearchMedicalPapersFlags.maxSources
	}
	if cmd.Flags().Changed("include-trials") {
		bodyMap["include_trials"] = datasetsSearchMedicalPapersFlags.includeTrials
	}
	if cmd.Flags().Changed("recency-years") {
		bodyMap["recency_years"] = datasetsSearchMedicalPapersFlags.recencyYears
	}
	if cmd.Flags().Changed("stream") {
		bodyMap["stream"] = datasetsSearchMedicalPapersFlags.stream
	}
	req.Body = bodyMap

	resp, err := client.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "deadline exceeded") {
			_invState.errorType = "timeout"
		} else {
			_invState.errorType = "network_error"
		}
		e := output.NetworkError(err)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if resp.StatusCode >= 400 {
		if resp.StatusCode >= 500 {
			_invState.errorType = "http_5xx"
		} else {
			_invState.errorType = "http_4xx"
		}
		_invState.errorCode = resp.StatusCode
		e := output.HTTPError(resp.StatusCode, resp.Body)
		e.Write(os.Stderr)
		return output.NewExitError(e)
	}

	if rootFlags.jq != "" {
		return output.JQFilter(_stdoutCounter, resp.Body, rootFlags.jq)
	}
	return output.Print(_stdoutCounter, resp.Body, output.Format(cfg.OutputFormat))
}
