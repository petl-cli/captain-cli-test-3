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

var companiesFinancialsCmd = &cobra.Command{
	Use:   "financials",
	Short: "Get Company Financials",
	RunE:  runCompaniesFinancials,
}

var companiesFinancialsFlags struct {
	companyId  string
	fiscalYear int
}

func init() {
	companiesFinancialsCmd.Flags().StringVar(&companiesFinancialsFlags.companyId, "company-id", "", "Company entity ID, website domain, or company name (e.g., 'openai.com', 'OpenAI', or UUID)")
	companiesFinancialsCmd.MarkFlagRequired("company-id")
	companiesFinancialsCmd.Flags().IntVar(&companiesFinancialsFlags.fiscalYear, "fiscal-year", 0, "Fiscal year (default: most recent)")

	companiesCmd.AddCommand(companiesFinancialsCmd)
}

func runCompaniesFinancials(cmd *cobra.Command, args []string) error {
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
			Name:        "fiscal-year",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Fiscal year (default: most recent)",
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

		schema := map[string]any{
			"command":     "financials",
			"description": "Get Company Financials",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/companies/{company_id}/financials",
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
	pathParams["company_id"] = fmt.Sprintf("%v", companiesFinancialsFlags.companyId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/companies/{company_id}/financials", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("fiscal-year") {
		req.QueryParams["fiscal_year"] = fmt.Sprintf("%v", companiesFinancialsFlags.fiscalYear)
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
