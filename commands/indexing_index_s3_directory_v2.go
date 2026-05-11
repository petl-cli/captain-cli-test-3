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

var indexingIndexS3DirectoryV2Cmd = &cobra.Command{
	Use:   "index-s3-directory-v2",
	Short: "Index S3 Directory",
	RunE:  runIndexingIndexS3DirectoryV2,
}

var indexingIndexS3DirectoryV2Flags struct {
	xOrganizationId    string
	collectionName     string
	idempotencyKey     string
	bucketName         string
	directoryPath      string
	bucketRegion       string
	awsAccessKeyId     string
	awsSecretAccessKey string
	processingType     string
	maxFiles           int
	skipExisting       bool
	parsingScript      string
	body               string
}

func init() {
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexS3DirectoryV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexS3DirectoryV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.idempotencyKey, "idempotency-key", "", "UUID for request deduplication")
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.bucketName, "bucket-name", "", "Name of the S3 bucket")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.directoryPath, "directory-path", "", "Path to the directory within the bucket. Accepts either a relative path (e.g., 'reports/2024/january') or a full S3 URI (e.g., 's3://my-bucket/reports/2024/january'). All files within this directory and its subdirectories will be indexed.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.bucketRegion, "bucket-region", "", "AWS region where the bucket is located")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.awsAccessKeyId, "aws-access-key-id", "", "AWS access key ID with read access to the bucket. Use this for long-lived IAM-user credentials. Omit when using the role-based 'auth' block.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.awsSecretAccessKey, "aws-secret-access-key", "", "AWS secret access key. Use this for long-lived IAM-user credentials. Omit when using the role-based 'auth' block.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.processingType, "processing-type", "", "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().IntVar(&indexingIndexS3DirectoryV2Flags.maxFiles, "max-files", 0, "Maximum number of files to index (optional)")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().BoolVar(&indexingIndexS3DirectoryV2Flags.skipExisting, "skip-existing", false, "Skip files that are already indexed in the collection. When true, only new files will be indexed. Set to false to re-index all files.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.parsingScript, "parsing-script", "", "Relative path to a JavaScript parsing script for JSON files (e.g. 'research/paper-parser'). When provided, .json files are processed through a sandboxed V8 isolate that executes the script to extract text and metadata. Without this parameter, .json files are indexed as raw text. Scripts are org-scoped and managed in the Parser Studio.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	indexingIndexS3DirectoryV2Cmd.Flags().StringVar(&indexingIndexS3DirectoryV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexS3DirectoryV2Cmd)
}

func runIndexingIndexS3DirectoryV2(cmd *cobra.Command, args []string) error {
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
			Name:        "bucket-name",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Name of the S3 bucket",
		})
		flags = append(flags, flagSchema{
			Name:        "directory-path",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Path to the directory within the bucket. Accepts either a relative path (e.g., 'reports/2024/january') or a full S3 URI (e.g., 's3://my-bucket/reports/2024/january'). All files within this directory and its subdirectories will be indexed.",
		})
		flags = append(flags, flagSchema{
			Name:        "bucket-region",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "AWS region where the bucket is located",
		})
		flags = append(flags, flagSchema{
			Name:        "aws-access-key-id",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "AWS access key ID with read access to the bucket. Use this for long-lived IAM-user credentials. Omit when using the role-based 'auth' block.",
		})
		flags = append(flags, flagSchema{
			Name:        "aws-secret-access-key",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "AWS secret access key. Use this for long-lived IAM-user credentials. Omit when using the role-based 'auth' block.",
		})
		flags = append(flags, flagSchema{
			Name:        "auth",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Cross-account role-assumption auth for S3 indexing. Captain calls sts:AssumeRole on the supplied role_arn (with the supplied external_id) instead of using static IAM-user keys, so no long-lived secrets cross the boundary.",
		})
		flags = append(flags, flagSchema{
			Name:        "processing-type",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "Document processing type. 'advanced' uses agentic OCR with AI-enhanced extraction for complex layouts, tables, figures, charts, and documents containing images. 'basic' provides reliable OCR optimized for general document indexing and high-volume processing.",
		})
		flags = append(flags, flagSchema{
			Name:        "max-files",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Maximum number of files to index (optional)",
		})
		flags = append(flags, flagSchema{
			Name:        "skip-existing",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Skip files that are already indexed in the collection. When true, only new files will be indexed. Set to false to re-index all files.",
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
			Description: "Indexing Job Started",
		})

		schema := map[string]any{
			"command":     "index-s3-directory-v2",
			"description": "Index S3 Directory",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/s3/directory",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexS3DirectoryV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/s3/directory", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexS3DirectoryV2Flags.xOrganizationId)
	}
	if cmd.Flags().Changed("idempotency-key") {
		req.Headers["Idempotency-Key"] = fmt.Sprintf("%v", indexingIndexS3DirectoryV2Flags.idempotencyKey)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexS3DirectoryV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexS3DirectoryV2Flags.body), &bodyMap); err != nil {
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
		bodyMap["bucket_name"] = indexingIndexS3DirectoryV2Flags.bucketName
	}
	if cmd.Flags().Changed("directory-path") {
		bodyMap["directory_path"] = indexingIndexS3DirectoryV2Flags.directoryPath
	}
	if cmd.Flags().Changed("bucket-region") {
		bodyMap["bucket_region"] = indexingIndexS3DirectoryV2Flags.bucketRegion
	}
	if cmd.Flags().Changed("aws-access-key-id") {
		bodyMap["aws_access_key_id"] = indexingIndexS3DirectoryV2Flags.awsAccessKeyId
	}
	if cmd.Flags().Changed("aws-secret-access-key") {
		bodyMap["aws_secret_access_key"] = indexingIndexS3DirectoryV2Flags.awsSecretAccessKey
	}
	if cmd.Flags().Changed("processing-type") {
		bodyMap["processing_type"] = indexingIndexS3DirectoryV2Flags.processingType
	}
	if cmd.Flags().Changed("max-files") {
		bodyMap["max_files"] = indexingIndexS3DirectoryV2Flags.maxFiles
	}
	if cmd.Flags().Changed("skip-existing") {
		bodyMap["skip_existing"] = indexingIndexS3DirectoryV2Flags.skipExisting
	}
	if cmd.Flags().Changed("parsing-script") {
		bodyMap["parsing_script"] = indexingIndexS3DirectoryV2Flags.parsingScript
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
