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

var limitedPartnersLpsSearchCmd = &cobra.Command{
	Use:   "lps-search",
	Short: "Search Limited Partners",
	RunE:  runLimitedPartnersLpsSearch,
}

var limitedPartnersLpsSearchFlags struct {
	q      string
	lpType string
	limit  int
}

func init() {
	limitedPartnersLpsSearchCmd.Flags().StringVar(&limitedPartnersLpsSearchFlags.q, "q", "", "LP name or keyword (e.g., 'CalPERS')")
	limitedPartnersLpsSearchCmd.MarkFlagRequired("q")
	limitedPartnersLpsSearchCmd.Flags().StringVar(&limitedPartnersLpsSearchFlags.lpType, "lp-type", "", "Filter by LP type (e.g., 'pension_fund', 'endowment', 'family_office', 'sovereign_wealth', 'insurance')")
	limitedPartnersLpsSearchCmd.Flags().IntVar(&limitedPartnersLpsSearchFlags.limit, "limit", 0, "Maximum number of results to return (1�100, default: 10)")

	limitedPartnersCmd.AddCommand(limitedPartnersLpsSearchCmd)
}

func runLimitedPartnersLpsSearch(cmd *cobra.Command, args []string) error {
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
			Description: "LP name or keyword (e.g., 'CalPERS')",
		})
		flags = append(flags, flagSchema{
			Name:        "lp-type",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by LP type (e.g., 'pension_fund', 'endowment', 'family_office', 'sovereign_wealth', 'insurance')",
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
			"command":     "lps-search",
			"description": "Search Limited Partners",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/limited-partners/search",
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
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/limited-partners/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", limitedPartnersLpsSearchFlags.q)
	}
	if cmd.Flags().Changed("lp-type") {
		req.QueryParams["lp_type"] = fmt.Sprintf("%v", limitedPartnersLpsSearchFlags.lpType)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", limitedPartnersLpsSearchFlags.limit)
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
