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

var indexingValidateParsingScriptV2Cmd = &cobra.Command{
	Use:   "validate-parsing-script-v2",
	Short: "Validate Parse Script",
	RunE:  runIndexingValidateParsingScriptV2,
}

var indexingValidateParsingScriptV2Flags struct {
	authorization   string
	xOrganizationId string
	file            string
	body            string
}

func init() {
	indexingValidateParsingScriptV2Cmd.Flags().StringVar(&indexingValidateParsingScriptV2Flags.authorization, "authorization", "", "Bearer token - your Captain API key.")
	indexingValidateParsingScriptV2Cmd.MarkFlagRequired("authorization")
	indexingValidateParsingScriptV2Cmd.Flags().StringVar(&indexingValidateParsingScriptV2Flags.xOrganizationId, "x-organization-id", "", "Organization UUID.")
	indexingValidateParsingScriptV2Cmd.MarkFlagRequired("x-organization-id")
	indexingValidateParsingScriptV2Cmd.Flags().StringVar(&indexingValidateParsingScriptV2Flags.file, "file", "", "The .js parsing script file to validate. Max 1 MB.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingValidateParsingScriptV2Cmd.Flags().StringVar(&indexingValidateParsingScriptV2Flags.body, "body", "", "Full request body as JSON (overrides individual flags)")

	indexingCmd.AddCommand(indexingValidateParsingScriptV2Cmd)
}

func runIndexingValidateParsingScriptV2(cmd *cobra.Command, args []string) error {
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
			Name:        "authorization",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "Bearer token - your Captain API key.",
		})
		flags = append(flags, flagSchema{
			Name:        "x-organization-id",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "Organization UUID.",
		})
		flags = append(flags, flagSchema{
			Name:        "file",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "The .js parsing script file to validate. Max 1 MB.",
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
			Description: "Validation result. Returned for both valid AND invalid scripts.",
		})

		schema := map[string]any{
			"command":     "validate-parsing-script-v2",
			"description": "Validate Parse Script",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/parsing-scripts/validate",
			},
			"input": map[string]any{
				"flags":         flags,
				"body_flag":     true,
				"body_required": true,
			},
			"output": map[string]any{
				"responses": responses,
			},
			"semantics": map[string]any{
				"safe":         false,
				"idempotent":   false,
				"reversible":   true,
				"side_effects": []string{"creates_resource"},
				"impact":       "medium",
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
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/parsing-scripts/validate", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("authorization") {
		req.Headers["Authorization"] = fmt.Sprintf("%v", indexingValidateParsingScriptV2Flags.authorization)
	}
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingValidateParsingScriptV2Flags.xOrganizationId)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingValidateParsingScriptV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingValidateParsingScriptV2Flags.body), &bodyMap); err != nil {
			_invState.errorType = "parse_error"
			cliErr := &output.CLIError{
				Error:    true,
				Code:     "validation_error",
				Message:  fmt.Sprintf("invalid JSON in --body: %v", err),
				ExitCode: output.ExitValidation,
			}
			cliErr.Write(os.Stderr)
			return output.NewExitError(cliErr)
		}
	}
	// Individual flags overlay onto body (flags take precedence over --body JSON)
	if cmd.Flags().Changed("file") {
		bodyMap["file"] = indexingValidateParsingScriptV2Flags.file
	}
	req.Body = bodyMap

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
