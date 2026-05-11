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

var datasetsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Dataset Articles",
	RunE:  runDatasetsSearch,
}

var datasetsSearchFlags struct {
	dataset string
	q       string
	limit   int
	author  string
}

func init() {
	datasetsSearchCmd.Flags().StringVar(&datasetsSearchFlags.dataset, "dataset", "", "The dataset to search. Contact your Account Executive for available datasets.")
	datasetsSearchCmd.MarkFlagRequired("dataset")
	datasetsSearchCmd.Flags().StringVar(&datasetsSearchFlags.q, "q", "", "Search query")
	datasetsSearchCmd.MarkFlagRequired("q")
	datasetsSearchCmd.Flags().IntVar(&datasetsSearchFlags.limit, "limit", 0, "Maximum number of results to return (default: 10, max: 100)")
	datasetsSearchCmd.Flags().StringVar(&datasetsSearchFlags.author, "author", "", "Filter results by author/byline name. Used as an AND condition with `q`  -  returns only articles matching BOTH the query topic AND the specified author. For all articles by an author regardless of topic, use a broad query like `q=*` with `author`.")

	datasetsCmd.AddCommand(datasetsSearchCmd)
}

func runDatasetsSearch(cmd *cobra.Command, args []string) error {
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
			Name:        "dataset",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "The dataset to search. Contact your Account Executive for available datasets.",
		})
		flags = append(flags, flagSchema{
			Name:        "q",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "Search query",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Maximum number of results to return (default: 10, max: 100)",
		})
		flags = append(flags, flagSchema{
			Name:        "author",
			Type:        "string",
			Required:    false,
			Location:    "query",
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
			Description: "Search completed successfully",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "",
			Description: "Invalid request - unknown dataset",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Unauthorized - invalid or missing API key",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "",
			Description: "Forbidden. Possible reasons: - API key does not belong to the specified organization - Trial expired and no payment method on file (`error: payment_required`) - Trial credits exhausted and no payment method on file (`error: credits_exhausted`)  For paid plans (Starter, Growth, Enterprise), requests are never blocked for credit usage  -  overage is billed at the end of the billing period.",
		})
		responses = append(responses, responseSchema{
			Status:      "503",
			ContentType: "",
			Description: "Search service unavailable",
		})

		schema := map[string]any{
			"command":     "search",
			"description": "Search Dataset Articles",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/{dataset}/search",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     false,
				"body_required": false,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         true,
				"idempotent":   true,
				"reversible":   true,
				"side_effects": []string{},
				"impact":       "low",
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
	pathParams["dataset"] = fmt.Sprintf("%v", datasetsSearchFlags.dataset)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/{dataset}/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", datasetsSearchFlags.q)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", datasetsSearchFlags.limit)
	}
	if cmd.Flags().Changed("author") {
		req.QueryParams["author"] = fmt.Sprintf("%v", datasetsSearchFlags.author)
	}

	// Header parameters

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
