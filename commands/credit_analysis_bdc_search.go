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

var creditAnalysisBdcSearchCmd = &cobra.Command{
	Use:   "bdc-search",
	Short: "Search BDC Portfolios",
	RunE:  runCreditAnalysisBdcSearch,
}

var creditAnalysisBdcSearchFlags struct {
	q              string
	seniority      string
	nonAccrualOnly bool
	limit          int
}

func init() {
	creditAnalysisBdcSearchCmd.Flags().StringVar(&creditAnalysisBdcSearchFlags.q, "q", "", "Search query  -  borrower name, industry, or keyword")
	creditAnalysisBdcSearchCmd.MarkFlagRequired("q")
	creditAnalysisBdcSearchCmd.Flags().StringVar(&creditAnalysisBdcSearchFlags.seniority, "seniority", "", "Filter by loan seniority")
	creditAnalysisBdcSearchCmd.Flags().BoolVar(&creditAnalysisBdcSearchFlags.nonAccrualOnly, "non-accrual-only", false, "Only return defaulted (non-accrual) investments")
	creditAnalysisBdcSearchCmd.Flags().IntVar(&creditAnalysisBdcSearchFlags.limit, "limit", 0, "Maximum results")

	creditAnalysisCmd.AddCommand(creditAnalysisBdcSearchCmd)
}

func runCreditAnalysisBdcSearch(cmd *cobra.Command, args []string) error {
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
			Required:    true,
			Location:    "query",
			Description: "Search query  -  borrower name, industry, or keyword",
		})
		flags = append(flags, flagSchema{
			Name:        "seniority",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by loan seniority",
		})
		flags = append(flags, flagSchema{
			Name:        "non-accrual-only",
			Type:        "boolean",
			Required:    false,
			Location:    "query",
			Description: "Only return defaulted (non-accrual) investments",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Maximum results",
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
			Description: "Matching BDC investments",
		})

		schema := map[string]any{
			"command":     "bdc-search",
			"description": "Search BDC Portfolios",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/credit-analysis/bdc/search",
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

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/credit-analysis/bdc/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", creditAnalysisBdcSearchFlags.q)
	}
	if cmd.Flags().Changed("seniority") {
		req.QueryParams["seniority"] = fmt.Sprintf("%v", creditAnalysisBdcSearchFlags.seniority)
	}
	if cmd.Flags().Changed("non-accrual-only") {
		req.QueryParams["non_accrual_only"] = fmt.Sprintf("%v", creditAnalysisBdcSearchFlags.nonAccrualOnly)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", creditAnalysisBdcSearchFlags.limit)
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
