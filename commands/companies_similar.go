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

var companiesSimilarCmd = &cobra.Command{
	Use:   "similar",
	Short: "Get Similar Companies",
	RunE:  runCompaniesSimilar,
}

var companiesSimilarFlags struct {
	companyId string
	limit     int
}

func init() {
	companiesSimilarCmd.Flags().StringVar(&companiesSimilarFlags.companyId, "company-id", "", "Company entity ID, website domain, or company name (e.g., 'openai.com', 'OpenAI', or UUID)")
	companiesSimilarCmd.MarkFlagRequired("company-id")
	companiesSimilarCmd.Flags().IntVar(&companiesSimilarFlags.limit, "limit", 0, "Maximum number of results to return (1�100, default: 10)")

	companiesCmd.AddCommand(companiesSimilarCmd)
}

func runCompaniesSimilar(cmd *cobra.Command, args []string) error {
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
			Name:        "company-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Company entity ID, website domain, or company name (e.g., 'openai.com', 'OpenAI', or UUID)",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Maximum number of results to return (1�100, default: 10)",
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
			Description: "Successful response",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Unauthorized - Invalid or missing API key",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "",
			Description: "Not found - Entity does not exist",
		})
		responses = append(responses, responseSchema{
			Status:      "501",
			ContentType: "",
			Description: "Not implemented - Data not available from public sources",
		})

		schema := map[string]any{
			"command":     "similar",
			"description": "Get Similar Companies",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/companies/{company_id}/similar",
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
	pathParams["company_id"] = fmt.Sprintf("%v", companiesSimilarFlags.companyId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/companies/{company_id}/similar", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", companiesSimilarFlags.limit)
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
