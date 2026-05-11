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

var collectionsDeleteDocumentV2Cmd = &cobra.Command{
	Use:   "delete-document-v2",
	Short: "Delete Document",
	RunE:  runCollectionsDeleteDocumentV2,
}

var collectionsDeleteDocumentV2Flags struct {
	collectionName string
	documentId     string
}

func init() {
	collectionsDeleteDocumentV2Cmd.Flags().StringVar(&collectionsDeleteDocumentV2Flags.collectionName, "collection-name", "", "Name of the collection")
	collectionsDeleteDocumentV2Cmd.MarkFlagRequired("collection-name")
	collectionsDeleteDocumentV2Cmd.Flags().StringVar(&collectionsDeleteDocumentV2Flags.documentId, "document-id", "", "ID of the document to delete")
	collectionsDeleteDocumentV2Cmd.MarkFlagRequired("document-id")

	collectionsCmd.AddCommand(collectionsDeleteDocumentV2Cmd)
}

func runCollectionsDeleteDocumentV2(cmd *cobra.Command, args []string) error {
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
			Name:        "collection-name",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Name of the collection",
		})
		flags = append(flags, flagSchema{
			Name:        "document-id",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "ID of the document to delete",
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
			Description: "Document Deleted",
		})
		responses = append(responses, responseSchema{
			Status:      "404",
			ContentType: "",
			Description: "Document not found",
		})

		schema := map[string]any{
			"command":     "delete-document-v2",
			"description": "Delete Document",
			"http": map[string]any{
				"method": "DELETE",
				"path":   "/v2/collections/{collection_name}/documents/{document_id}",
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
				"safe":         false,
				"idempotent":   true,
				"reversible":   false,
				"side_effects": []string{"destroys_resource"},
				"impact":       "high",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", collectionsDeleteDocumentV2Flags.collectionName)
	pathParams["document_id"] = fmt.Sprintf("%v", collectionsDeleteDocumentV2Flags.documentId)

	req := &httpclient.Request{
		Method:      "DELETE",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/documents/{document_id}", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

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
