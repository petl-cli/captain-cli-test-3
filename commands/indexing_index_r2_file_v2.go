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

var indexingIndexR2FileV2Cmd = &cobra.Command{
	Use:   "index-r2-file-v2",
	Short: "Index R2 File",
	RunE:  runIndexingIndexR2FileV2,
}

var indexingIndexR2FileV2Flags struct {
	xOrganizationId string
	collectionName  string
	bucketName      string
	fileUri         string
	accountId       string
	accessKeyId     string
	secretAccessKey string
	jurisdiction    string
	processingType  string
	parsingScript   string
	body            string
}

func init() {
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexR2FileV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexR2FileV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.bucketName, "bucket-name", "", "Name of the R2 bucket")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.fileUri, "file-uri", "", "R2 URI format: r2://bucket-name/path/to/file.pdf")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.accountId, "account-id", "", "Cloudflare account ID (found in your R2 dashboard URL)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.accessKeyId, "access-key-id", "", "R2 S3 API token Access Key ID")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.secretAccessKey, "secret-access-key", "", "R2 S3 API token Secret Access Key")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.jurisdiction, "jurisdiction", "", "R2 jurisdiction. 'default' for global, 'eu' for EU-only storage, 'fedramp' for FedRAMP-compliant storage.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.processingType, "processing-type", "", "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.parsingScript, "parsing-script", "", "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexR2FileV2Cmd.Flags().StringVar(&indexingIndexR2FileV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexR2FileV2Cmd)
}

func runIndexingIndexR2FileV2(cmd *cobra.Command, args []string) error {
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
			Name:        "bucket-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the R2 bucket",
		})
		flags = append(flags, flagSchema{
			Name:        "file-uri",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "R2 URI format: r2://bucket-name/path/to/file.pdf",
		})
		flags = append(flags, flagSchema{
			Name:        "account-id",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Cloudflare account ID (found in your R2 dashboard URL)",
		})
		flags = append(flags, flagSchema{
			Name:        "access-key-id",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "R2 S3 API token Access Key ID",
		})
		flags = append(flags, flagSchema{
			Name:        "secret-access-key",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "R2 S3 API token Secret Access Key",
		})
		flags = append(flags, flagSchema{
			Name:        "jurisdiction",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "R2 jurisdiction. 'default' for global, 'eu' for EU-only storage, 'fedramp' for FedRAMP-compliant storage.",
		})
		flags = append(flags, flagSchema{
			Name:        "processing-type",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.",
		})
		flags = append(flags, flagSchema{
			Name:        "custom-metadata",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Custom metadata to attach to all chunks from this file. Keys must be strings. Values: str, int, float, bool, or array of strings.",
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
			Description: "Indexing Job Started",
		})

		schema := map[string]any{
			"command":     "index-r2-file-v2",
			"description": "Index R2 File",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/r2/file",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexR2FileV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/r2/file", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexR2FileV2Flags.xOrganizationId)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexR2FileV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexR2FileV2Flags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("bucket-name") {
		bodyMap["bucket_name"] = indexingIndexR2FileV2Flags.bucketName
	}
	if cmd.Flags().Changed("file-uri") {
		bodyMap["file_uri"] = indexingIndexR2FileV2Flags.fileUri
	}
	if cmd.Flags().Changed("account-id") {
		bodyMap["account_id"] = indexingIndexR2FileV2Flags.accountId
	}
	if cmd.Flags().Changed("access-key-id") {
		bodyMap["access_key_id"] = indexingIndexR2FileV2Flags.accessKeyId
	}
	if cmd.Flags().Changed("secret-access-key") {
		bodyMap["secret_access_key"] = indexingIndexR2FileV2Flags.secretAccessKey
	}
	if cmd.Flags().Changed("jurisdiction") {
		bodyMap["jurisdiction"] = indexingIndexR2FileV2Flags.jurisdiction
	}
	if cmd.Flags().Changed("processing-type") {
		bodyMap["processing_type"] = indexingIndexR2FileV2Flags.processingType
	}
	if cmd.Flags().Changed("parsing-script") {
		bodyMap["parsing_script"] = indexingIndexR2FileV2Flags.parsingScript
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
