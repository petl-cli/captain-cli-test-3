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

var indexingIndexUrlV2Cmd = &cobra.Command{
	Use:   "index-url-v2",
	Short: "Index URLs",
	RunE:  runIndexingIndexUrlV2,
}

var indexingIndexUrlV2Flags struct {
	xOrganizationId string
	collectionName  string
	idempotencyKey  string
	url             string
	urls            []string
	processingType  string
	parsingScript   string
	body            string
}

func init() {
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexUrlV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexUrlV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.idempotencyKey, "idempotency-key", "", "UUID for request deduplication")
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.url, "url", "", "A single public URL to a document or web page. Hosted files (PDF, DOCX, etc.) are indexed directly. Web pages (HTML) are automatically scraped  -  text and images are extracted. Provide either 'url' or 'urls', not both.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexUrlV2Cmd.Flags().StringSliceVar(&indexingIndexUrlV2Flags.urls, "urls", nil, "An array of public URLs to documents or web pages. Each URL is auto-detected  -  hosted files are indexed directly, web pages are scraped. Provide either 'url' or 'urls', not both.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.processingType, "processing-type", "", "Processing mode. For hosted documents: 'advanced' enables AI-enhanced extraction for complex layouts, tables, figures, and charts; 'basic' provides standard document processing. For web pages: 'advanced' extracts both text content and page images; 'basic' extracts text content only (faster, lower cost).")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.parsingScript, "parsing-script", "", "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexUrlV2Cmd.Flags().StringVar(&indexingIndexUrlV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexUrlV2Cmd)
}

func runIndexingIndexUrlV2(cmd *cobra.Command, args []string) error {
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
			Name:        "x-organization-id",
			Type:        "string",
			Required:    true,
			Location:    "header",
			Description: "The organization ID to scope the request",
		})
		flags = append(flags, flagSchema{
			Name:        "collection-name",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Name of the collection to index into",
		})
		flags = append(flags, flagSchema{
			Name:        "idempotency-key",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "UUID for request deduplication",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "A single public URL to a document or web page. Hosted files (PDF, DOCX, etc.) are indexed directly. Web pages (HTML) are automatically scraped  -  text and images are extracted. Provide either 'url' or 'urls', not both.",
		})
		flags = append(flags, flagSchema{
			Name:        "urls",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "An array of public URLs to documents or web pages. Each URL is auto-detected  -  hosted files are indexed directly, web pages are scraped. Provide either 'url' or 'urls', not both.",
		})
		flags = append(flags, flagSchema{
			Name:        "processing-type",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Processing mode. For hosted documents: 'advanced' enables AI-enhanced extraction for complex layouts, tables, figures, and charts; 'basic' provides standard document processing. For web pages: 'advanced' extracts both text content and page images; 'basic' extracts text content only (faster, lower cost).",
		})
		flags = append(flags, flagSchema{
			Name:        "custom-metadata",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Custom metadata to attach to all indexed chunks. Keys must be strings. Values: str, int, float, bool, or array of strings.",
		})
		flags = append(flags, flagSchema{
			Name:        "parsing-script",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.",
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
			Description: "Indexing job started",
		})

		schema := map[string]any{
			"command":     "index-url-v2",
			"description": "Index URLs",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/url",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexUrlV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/url", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexUrlV2Flags.xOrganizationId)
	}
	if cmd.Flags().Changed("idempotency-key") {
		req.Headers["Idempotency-Key"] = fmt.Sprintf("%v", indexingIndexUrlV2Flags.idempotencyKey)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexUrlV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexUrlV2Flags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("url") {
		bodyMap["url"] = indexingIndexUrlV2Flags.url
	}
	if cmd.Flags().Changed("urls") {
		bodyMap["urls"] = indexingIndexUrlV2Flags.urls
	}
	if cmd.Flags().Changed("processing-type") {
		bodyMap["processing_type"] = indexingIndexUrlV2Flags.processingType
	}
	if cmd.Flags().Changed("parsing-script") {
		bodyMap["parsing_script"] = indexingIndexUrlV2Flags.parsingScript
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
