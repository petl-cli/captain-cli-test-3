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

var indexingIndexAzureFileV2Cmd = &cobra.Command{
	Use:   "index-azure-file-v2",
	Short: "Index Azure File",
	RunE:  runIndexingIndexAzureFileV2,
}

var indexingIndexAzureFileV2Flags struct {
	xOrganizationId string
	collectionName  string
	containerName   string
	fileUri         string
	accountName     string
	accountKey      string
	processingType  string
	parsingScript   string
	body            string
}

func init() {
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexAzureFileV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexAzureFileV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.containerName, "container-name", "", "Name of the Azure Blob Storage container")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.fileUri, "file-uri", "", "Azure Blob Storage URI format: https://{account}.blob.core.windows.net/{container}/path/to/file.pdf")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.accountName, "account-name", "", "Azure Storage account name")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.accountKey, "account-key", "", "Azure Storage account key (base64-encoded)")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.processingType, "processing-type", "", "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.parsingScript, "parsing-script", "", "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexAzureFileV2Cmd.Flags().StringVar(&indexingIndexAzureFileV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexAzureFileV2Cmd)
}

func runIndexingIndexAzureFileV2(cmd *cobra.Command, args []string) error {
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
			Name:        "container-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the Azure Blob Storage container",
		})
		flags = append(flags, flagSchema{
			Name:        "file-uri",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Azure Blob Storage URI format: https://{account}.blob.core.windows.net/{container}/path/to/file.pdf",
		})
		flags = append(flags, flagSchema{
			Name:        "account-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Azure Storage account name",
		})
		flags = append(flags, flagSchema{
			Name:        "account-key",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Azure Storage account key (base64-encoded)",
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
			"command":     "index-azure-file-v2",
			"description": "Index Azure File",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/azure/file",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexAzureFileV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/azure/file", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexAzureFileV2Flags.xOrganizationId)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexAzureFileV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexAzureFileV2Flags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("container-name") {
		bodyMap["container_name"] = indexingIndexAzureFileV2Flags.containerName
	}
	if cmd.Flags().Changed("file-uri") {
		bodyMap["file_uri"] = indexingIndexAzureFileV2Flags.fileUri
	}
	if cmd.Flags().Changed("account-name") {
		bodyMap["account_name"] = indexingIndexAzureFileV2Flags.accountName
	}
	if cmd.Flags().Changed("account-key") {
		bodyMap["account_key"] = indexingIndexAzureFileV2Flags.accountKey
	}
	if cmd.Flags().Changed("processing-type") {
		bodyMap["processing_type"] = indexingIndexAzureFileV2Flags.processingType
	}
	if cmd.Flags().Changed("parsing-script") {
		bodyMap["parsing_script"] = indexingIndexAzureFileV2Flags.parsingScript
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
