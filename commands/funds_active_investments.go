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

var fundsActiveInvestmentsCmd = &cobra.Command{
	Use:   "active-investments",
	Short: "Get Fund Active Investments",
	RunE:  runFundsActiveInvestments,
}

var fundsActiveInvestmentsFlags struct {
	fundId   string
	page     int
	pageSize int
}

func init() {
	fundsActiveInvestmentsCmd.Flags().StringVar(&fundsActiveInvestmentsFlags.fundId, "fund-id", "", "Fund entity ID")
	fundsActiveInvestmentsCmd.MarkFlagRequired("fund-id")
	fundsActiveInvestmentsCmd.Flags().IntVar(&fundsActiveInvestmentsFlags.page, "page", 0, "Page number for pagination (default: 1)")
	fundsActiveInvestmentsCmd.Flags().IntVar(&fundsActiveInvestmentsFlags.pageSize, "page-size", 0, "Results per page (default: 50, max: 1000)")

	fundsCmd.AddCommand(fundsActiveInvestmentsCmd)
}

func runFundsActiveInvestments(cmd *cobra.Command, args []string) error {
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
			Name:        "fund-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Fund entity ID",
		})
		flags = append(flags, flagSchema{
			Name:        "page",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Page number for pagination (default: 1)",
		})
		flags = append(flags, flagSchema{
			Name:        "page-size",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Results per page (default: 50, max: 1000)",
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
			"command":     "active-investments",
			"description": "Get Fund Active Investments",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/funds/{fund_id}/active-investments",
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
	pathParams["fund_id"] = fmt.Sprintf("%v", fundsActiveInvestmentsFlags.fundId)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/funds/{fund_id}/active-investments", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("page") {
		req.QueryParams["page"] = fmt.Sprintf("%v", fundsActiveInvestmentsFlags.page)
	}
	if cmd.Flags().Changed("page-size") {
		req.QueryParams["page_size"] = fmt.Sprintf("%v", fundsActiveInvestmentsFlags.pageSize)
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
