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

var dealsSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search Deals",
	RunE:  runDealsSearch,
}

var dealsSearchFlags struct {
	q         string
	company   string
	dealType  string
	minAmount float64
	maxAmount float64
	startDate string
	endDate   string
	page      int
	pageSize  int
}

func init() {
	dealsSearchCmd.Flags().StringVar(&dealsSearchFlags.q, "q", "", "Company name or deal keyword (e.g., 'OpenAI Series B')")
	dealsSearchCmd.Flags().StringVar(&dealsSearchFlags.company, "company", "", "Filter by company name or domain (e.g., 'OpenAI')")
	dealsSearchCmd.Flags().StringVar(&dealsSearchFlags.dealType, "deal-type", "", "Filter by deal type (e.g., 'series_a', 'series_b', 'seed', 'ipo', 'acquisition', 'debt')")
	dealsSearchCmd.Flags().Float64Var(&dealsSearchFlags.minAmount, "min-amount", 0, "Minimum deal amount")
	dealsSearchCmd.Flags().Float64Var(&dealsSearchFlags.maxAmount, "max-amount", 0, "Maximum deal amount")
	dealsSearchCmd.Flags().StringVar(&dealsSearchFlags.startDate, "start-date", "", "Start date (YYYY-MM-DD)")
	dealsSearchCmd.Flags().StringVar(&dealsSearchFlags.endDate, "end-date", "", "End date (YYYY-MM-DD)")
	dealsSearchCmd.Flags().IntVar(&dealsSearchFlags.page, "page", 0, "Page number")
	dealsSearchCmd.Flags().IntVar(&dealsSearchFlags.pageSize, "page-size", 0, "Results per page")

	dealsCmd.AddCommand(dealsSearchCmd)
}

func runDealsSearch(cmd *cobra.Command, args []string) error {
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
			Location:    "query",
			Description: "Company name or deal keyword (e.g., 'OpenAI Series B')",
		})
		flags = append(flags, flagSchema{
			Name:        "company",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by company name or domain (e.g., 'OpenAI')",
		})
		flags = append(flags, flagSchema{
			Name:        "deal-type",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by deal type (e.g., 'series_a', 'series_b', 'seed', 'ipo', 'acquisition', 'debt')",
		})
		flags = append(flags, flagSchema{
			Name:        "min-amount",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Minimum deal amount",
		})
		flags = append(flags, flagSchema{
			Name:        "max-amount",
			Type:        "number",
			Required:    false,
			Location:    "query",
			Description: "Maximum deal amount",
		})
		flags = append(flags, flagSchema{
			Name:        "start-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Start date (YYYY-MM-DD)",
		})
		flags = append(flags, flagSchema{
			Name:        "end-date",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "End date (YYYY-MM-DD)",
		})
		flags = append(flags, flagSchema{
			Name:        "page",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Page number",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Results per page",
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
			"description": "Search Deals",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/deals/search",
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
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/deals/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", dealsSearchFlags.q)
	}
	if cmd.Flags().Changed("company") {
		req.QueryParams["company"] = fmt.Sprintf("%v", dealsSearchFlags.company)
	}
	if cmd.Flags().Changed("deal-type") {
		req.QueryParams["deal_type"] = fmt.Sprintf("%v", dealsSearchFlags.dealType)
	}
	if cmd.Flags().Changed("min-amount") {
		req.QueryParams["min_amount"] = fmt.Sprintf("%v", dealsSearchFlags.minAmount)
	}
	if cmd.Flags().Changed("max-amount") {
		req.QueryParams["max_amount"] = fmt.Sprintf("%v", dealsSearchFlags.maxAmount)
	}
	if cmd.Flags().Changed("start-date") {
		req.QueryParams["start_date"] = fmt.Sprintf("%v", dealsSearchFlags.startDate)
	}
	if cmd.Flags().Changed("end-date") {
		req.QueryParams["end_date"] = fmt.Sprintf("%v", dealsSearchFlags.endDate)
	}
	if cmd.Flags().Changed("page") {
		req.QueryParams["page"] = fmt.Sprintf("%v", dealsSearchFlags.page)
	}
	if cmd.Flags().Changed("page-size") {
		req.QueryParams["page_size"] = fmt.Sprintf("%v", dealsSearchFlags.pageSize)
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
