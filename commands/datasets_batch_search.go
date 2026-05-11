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

var datasetsBatchSearchCmd = &cobra.Command{
	Use:   "batch-search",
	Short: "Batch Search Articles",
	RunE:  runDatasetsBatchSearch,
}

var datasetsBatchSearchFlags struct {
	q        string
	datasets []string
	limit    int
	author   string
	body     string
}

func init() {
	datasetsBatchSearchCmd.Flags().StringVar(&datasetsBatchSearchFlags.q, "q", "", "Search query")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsBatchSearchCmd.Flags().StringSliceVar(&datasetsBatchSearchFlags.datasets, "datasets", nil, "List of dataset names to search. Defaults to all datasets if not provided.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsBatchSearchCmd.Flags().IntVar(&datasetsBatchSearchFlags.limit, "limit", 0, "Maximum number of results to return (default: 10, max: 100)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsBatchSearchCmd.Flags().StringVar(&datasetsBatchSearchFlags.author, "author", "", "Filter results by author/byline name. Used as an AND condition with `q`  -  returns only articles matching BOTH the query topic AND the specified author. For all articles by an author regardless of topic, use a broad query like `q=*` with `author`.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	datasetsBatchSearchCmd.Flags().StringVar(&datasetsBatchSearchFlags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	datasetsCmd.AddCommand(datasetsBatchSearchCmd)
}

func runDatasetsBatchSearch(cmd *cobra.Command, args []string) error {
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
			Name:        "q",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Search query",
		})
		flags = append(flags, flagSchema{
			Name:        "datasets",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "List of dataset names to search. Defaults to all datasets if not provided.",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Maximum number of results to return (default: 10, max: 100)",
		})
		flags = append(flags, flagSchema{
			Name:        "author",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Filter results by author/byline name. Used as an AND condition with `q`  -  returns only articles matching BOTH the query topic AND the specified author. For all articles by an author regardless of topic, use a broad query like `q=*` with `author`.",
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
			Description: "Batch search completed successfully",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "",
			Description: "Invalid request - unknown dataset(s)",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Unauthorized - invalid or missing API key",
		})
		responses = append(responses, responseSchema{
			Status:      "503",
			ContentType: "",
			Description: "Search service unavailable",
		})

		schema := map[string]any{
			"command":     "batch-search",
			"description": "Batch Search Articles",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/datasets/batch-search",
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
		Path:        httpclient.SubstitutePath("/v2/datasets/batch-search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters

	// Request body
	bodyMap := map[string]any{}
	if datasetsBatchSearchFlags.body != "" {
		if err := json.Unmarshal([]byte(datasetsBatchSearchFlags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("q") {
		bodyMap["q"] = datasetsBatchSearchFlags.q
	}
	if cmd.Flags().Changed("datasets") {
		bodyMap["datasets"] = datasetsBatchSearchFlags.datasets
	}
	if cmd.Flags().Changed("limit") {
		bodyMap["limit"] = datasetsBatchSearchFlags.limit
	}
	if cmd.Flags().Changed("author") {
		bodyMap["author"] = datasetsBatchSearchFlags.author
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
