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

var fundsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Funds",
	RunE:  runFundsSearch,
}

var fundsSearchFlags struct {
	q           string
	fundType    string
	vintageYear int
	limit       int
}

func init() {
	fundsSearchCmd.Flags().StringVar(&fundsSearchFlags.q, "q", "", "Fund name or keyword (e.g., 'Sequoia Capital Fund XIV')")
	fundsSearchCmd.MarkFlagRequired("q")
	fundsSearchCmd.Flags().StringVar(&fundsSearchFlags.fundType, "fund-type", "", "Filter by fund type (e.g., 'venture', 'buyout', 'growth_equity', 'real_estate', 'debt')")
	fundsSearchCmd.Flags().IntVar(&fundsSearchFlags.vintageYear, "vintage-year", 0, "Filter by fund vintage year (e.g., 2023)")
	fundsSearchCmd.Flags().IntVar(&fundsSearchFlags.limit, "limit", 0, "Maximum number of results to return (1�100, default: 10)")

	fundsCmd.AddCommand(fundsSearchCmd)
}

func runFundsSearch(cmd *cobra.Command, args []string) error {
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
			Description: "Fund name or keyword (e.g., 'Sequoia Capital Fund XIV')",
		})
		flags = append(flags, flagSchema{
			Name:        "fund-type",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by fund type (e.g., 'venture', 'buyout', 'growth_equity', 'real_estate', 'debt')",
		})
		flags = append(flags, flagSchema{
			Name:        "vintage-year",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Filter by fund vintage year (e.g., 2023)",
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

		schema := map[string]any{
			"command":     "search",
			"description": "Search Funds",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/funds/search",
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
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/funds/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", fundsSearchFlags.q)
	}
	if cmd.Flags().Changed("fund-type") {
		req.QueryParams["fund_type"] = fmt.Sprintf("%v", fundsSearchFlags.fundType)
	}
	if cmd.Flags().Changed("vintage-year") {
		req.QueryParams["vintage_year"] = fmt.Sprintf("%v", fundsSearchFlags.vintageYear)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", fundsSearchFlags.limit)
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
