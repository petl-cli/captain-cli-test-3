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

var creditAnalysisBdcPortfolioCmd = &cobra.Command{
	Use:   "bdc-portfolio",
	Short: "Get BDC Portfolio",
	RunE:  runCreditAnalysisBdcPortfolio,
}

var creditAnalysisBdcPortfolioFlags struct {
	ticker  string
	quarter string
	limit   int
}

func init() {
	creditAnalysisBdcPortfolioCmd.Flags().StringVar(&creditAnalysisBdcPortfolioFlags.ticker, "ticker", "", "BDC ticker symbol (e.g., ARCC, BXSL, FSK)")
	creditAnalysisBdcPortfolioCmd.MarkFlagRequired("ticker")
	creditAnalysisBdcPortfolioCmd.Flags().StringVar(&creditAnalysisBdcPortfolioFlags.quarter, "quarter", "", "Quarter like '2025-Q1'  -  defaults to latest filing")
	creditAnalysisBdcPortfolioCmd.Flags().IntVar(&creditAnalysisBdcPortfolioFlags.limit, "limit", 0, "Maximum investments to return")

	creditAnalysisCmd.AddCommand(creditAnalysisBdcPortfolioCmd)
}

func runCreditAnalysisBdcPortfolio(cmd *cobra.Command, args []string) error {
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
			Name:        "ticker",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "BDC ticker symbol (e.g., ARCC, BXSL, FSK)",
		})
		flags = append(flags, flagSchema{
			Name:        "quarter",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Quarter like '2025-Q1'  -  defaults to latest filing",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Maximum investments to return",
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
			Description: "BDC portfolio holdings",
		})

		schema := map[string]any{
			"command":     "bdc-portfolio",
			"description": "Get BDC Portfolio",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/credit-analysis/bdc/{ticker}/portfolio",
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
	pathParams["ticker"] = fmt.Sprintf("%v", creditAnalysisBdcPortfolioFlags.ticker)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/credit-analysis/bdc/{ticker}/portfolio", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("quarter") {
		req.QueryParams["quarter"] = fmt.Sprintf("%v", creditAnalysisBdcPortfolioFlags.quarter)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", creditAnalysisBdcPortfolioFlags.limit)
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
