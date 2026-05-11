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

var indexingIndexGcsFileV2Cmd = &cobra.Command{
	Use:   "index-gcs-file-v2",
	Short: "Index GCS File",
	RunE:  runIndexingIndexGcsFileV2,
}

var indexingIndexGcsFileV2Flags struct {
	xOrganizationId    string
	collectionName     string
	bucketName         string
	fileUri            string
	serviceAccountJson string
	processingType     string
	parsingScript      string
	body               string
}

func init() {
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexGcsFileV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexGcsFileV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.bucketName, "bucket-name", "", "Name of the GCS bucket")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.fileUri, "file-uri", "", "GCS URI format: gs://bucket-name/path/to/file.pdf")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.serviceAccountJson, "service-account-json", "", "GCP service account JSON key with read access to the bucket")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.processingType, "processing-type", "", "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.parsingScript, "parsing-script", "", "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexGcsFileV2Cmd.Flags().StringVar(&indexingIndexGcsFileV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexGcsFileV2Cmd)
}

func runIndexingIndexGcsFileV2(cmd *cobra.Command, args []string) error {
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
			Description: "Name of the GCS bucket",
		})
		flags = append(flags, flagSchema{
			Name:        "file-uri",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "GCS URI format: gs://bucket-name/path/to/file.pdf",
		})
		flags = append(flags, flagSchema{
			Name:        "service-account-json",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "GCP service account JSON key with read access to the bucket",
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
			"command":     "index-gcs-file-v2",
			"description": "Index GCS File",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/gcs/file",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexGcsFileV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/gcs/file", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexGcsFileV2Flags.xOrganizationId)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexGcsFileV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexGcsFileV2Flags.body), &bodyMap); err != nil {
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
		bodyMap["bucket_name"] = indexingIndexGcsFileV2Flags.bucketName
	}
	if cmd.Flags().Changed("file-uri") {
		bodyMap["file_uri"] = indexingIndexGcsFileV2Flags.fileUri
	}
	if cmd.Flags().Changed("service-account-json") {
		bodyMap["service_account_json"] = indexingIndexGcsFileV2Flags.serviceAccountJson
	}
	if cmd.Flags().Changed("processing-type") {
		bodyMap["processing_type"] = indexingIndexGcsFileV2Flags.processingType
	}
	if cmd.Flags().Changed("parsing-script") {
		bodyMap["parsing_script"] = indexingIndexGcsFileV2Flags.parsingScript
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
