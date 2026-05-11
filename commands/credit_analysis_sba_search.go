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

var creditAnalysisSbaSearchCmd = &cobra.Command{
	Use:   "sba-search",
	Short: "Search SBA Loans",
	RunE:  runCreditAnalysisSbaSearch,
}

var creditAnalysisSbaSearchFlags struct {
	borrower  string
	lender    string
	state     string
	naics     string
	status    string
	minAmount float64
	limit     int
}

func init() {
	creditAnalysisSbaSearchCmd.Flags().StringVar(&creditAnalysisSbaSearchFlags.borrower, "borrower", "", "Borrower name (partial match)")
	creditAnalysisSbaSearchCmd.Flags().StringVar(&creditAnalysisSbaSearchFlags.lender, "lender", "", "Bank/lender name (partial match)")
	creditAnalysisSbaSearchCmd.Flags().StringVar(&creditAnalysisSbaSearchFlags.state, "state", "", "State code (e.g., CA, NY, TX)")
	creditAnalysisSbaSearchCmd.Flags().StringVar(&creditAnalysisSbaSearchFlags.naics, "naics", "", "NAICS code prefix (e.g., 5112 for software)")
	creditAnalysisSbaSearchCmd.Flags().StringVar(&creditAnalysisSbaSearchFlags.status, "status", "", "Loan status: CHGOFF (charged off), PIF (paid in full)")
	creditAnalysisSbaSearchCmd.Flags().Float64Var(&creditAnalysisSbaSearchFlags.minAmount, "min-amount", 0, "Minimum gross approval amount ($)")
	creditAnalysisSbaSearchCmd.Flags().IntVar(&creditAnalysisSbaSearchFlags.limit, "limit", 0, "Maximum results")

	creditAnalysisCmd.AddCommand(creditAnalysisSbaSearchCmd)
}

func runCreditAnalysisSbaSearch(cmd *cobra.Command, args []string) error {
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
			Name:        "borrower",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Borrower name (partial match)",
		})
		flags = append(flags, flagSchema{
			Name:        "lender",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Bank/lender name (partial match)",
		})
		flags = append(flags, flagSchema{
			Name:        "state",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "State code (e.g., CA, NY, TX)",
		})
		flags = append(flags, flagSchema{
			Name:        "naics",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "NAICS code prefix (e.g., 5112 for software)",
		})
		flags = append(flags, flagSchema{
			Name:        "status",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Loan status: CHGOFF (charged off), PIF (paid in full)",
		})
		flags = append(flags, flagSchema{
			Name:        "min-amount",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Minimum gross approval amount ($)",
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
			Description: "Matching SBA loans",
		})

		schema := map[string]any{
			"command":     "sba-search",
			"description": "Search SBA Loans",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/credit-analysis/sba/search",
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
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/credit-analysis/sba/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("borrower") {
		req.QueryParams["borrower"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.borrower)
	}
	if cmd.Flags().Changed("lender") {
		req.QueryParams["lender"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.lender)
	}
	if cmd.Flags().Changed("state") {
		req.QueryParams["state"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.state)
	}
	if cmd.Flags().Changed("naics") {
		req.QueryParams["naics"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.naics)
	}
	if cmd.Flags().Changed("status") {
		req.QueryParams["status"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.status)
	}
	if cmd.Flags().Changed("min-amount") {
		req.QueryParams["min_amount"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.minAmount)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", creditAnalysisSbaSearchFlags.limit)
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
