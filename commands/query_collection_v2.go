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

var queryCollectionV2Cmd = &cobra.Command{
	Use:   "collection-v2",
	Short: "Query Collection",
	RunE:  runQueryCollectionV2,
}

var queryCollectionV2Flags struct {
	xOrganizationId string
	collectionName  string
	idempotencyKey  string
	query           string
	inference       bool
	stream          bool
	topK            int
	rerank          bool
	customPrompt    string
	includeBbox     bool
	searchResults   bool
	body            string
}

func init() {
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	queryCollectionV2Cmd.MarkFlagRequired("x-organization-id")
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.collectionName, "collection-name", "", "Name of the collection to query")
	queryCollectionV2Cmd.MarkFlagRequired("collection-name")
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.idempotencyKey, "idempotency-key", "", "UUID for request deduplication")
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.query, "query", "", "The natural language query to search for")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().BoolVar(&queryCollectionV2Flags.inference, "inference", false, "Enable LLM-generated answers based on the relevant sections retrieved. When false, returns raw search results.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().BoolVar(&queryCollectionV2Flags.stream, "stream", false, "Enable real-time streaming of the response")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().IntVar(&queryCollectionV2Flags.topK, "top-k", 0, "Number of results to return. Only valid when inference=false. Not supported when inference=true (the agent controls its own search strategy).")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().BoolVar(&queryCollectionV2Flags.rerank, "rerank", false, "Enable reranking for improved relevance ordering. Uses Gemini Flash 2.5 by default, or Voyage AI rerank-2.5 as fallback. Adds ~100-300ms latency.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.customPrompt, "custom-prompt", "", "Custom system prompt to override the default RAG prompt when inference=true. Allows customizing how the LLM processes and responds to the query with the retrieved context.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().BoolVar(&queryCollectionV2Flags.includeBbox, "include-bbox", false, "Include normalized bounding box layout data for each search result. Returns element-level positions (titles, paragraphs, tables, figures, form fields) with page coordinates for PDF and DOCX files. Only supported with inference=false.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().BoolVar(&queryCollectionV2Flags.searchResults, "search-results", false, "When inference=true, include the raw search result chunks that were used as context for the LLM response. Defaults to false. Always true when inference=false.")
	// Note: body fields are not MarkFlagRequired in JSON mode — --body satisfies them too.
	queryCollectionV2Cmd.Flags().StringVar(&queryCollectionV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	queryCmd.AddCommand(queryCollectionV2Cmd)
}

func runQueryCollectionV2(cmd *cobra.Command, args []string) error {
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
			Description: "Name of the collection to query",
		})
		flags = append(flags, flagSchema{
			Name:        "idempotency-key",
			Type:        "string",
			Required:    false,
			Location:    "header",
			Description: "UUID for request deduplication",
		})
		flags = append(flags, flagSchema{
			Name:        "query",
			Type:        "string",
			Required:    true,
			Location:    "body",
			Description: "The natural language query to search for",
		})
		flags = append(flags, flagSchema{
			Name:        "inference",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Enable LLM-generated answers based on the relevant sections retrieved. When false, returns raw search results.",
		})
		flags = append(flags, flagSchema{
			Name:        "stream",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Enable real-time streaming of the response",
		})
		flags = append(flags, flagSchema{
			Name:        "top-k",
			Type:        "integer",
			Required:    false,
			Location:    "body",
			Description: "Number of results to return. Only valid when inference=false. Not supported when inference=true (the agent controls its own search strategy).",
		})
		flags = append(flags, flagSchema{
			Name:        "rerank",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Enable reranking for improved relevance ordering. Uses Gemini Flash 2.5 by default, or Voyage AI rerank-2.5 as fallback. Adds ~100-300ms latency.",
		})
		flags = append(flags, flagSchema{
			Name:        "metadata-filter",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Filter expression for vector search. Supports: $eq, $ne, $gt, $gte, $lt, $lte, $in, $nin, $and, $or",
		})
		flags = append(flags, flagSchema{
			Name:        "custom-prompt",
			Type:        "string",
			Required:    false,
			Location:    "body",
			Description: "Custom system prompt to override the default RAG prompt when inference=true. Allows customizing how the LLM processes and responds to the query with the retrieved context.",
		})
		flags = append(flags, flagSchema{
			Name:        "include-bbox",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "Include normalized bounding box layout data for each search result. Returns element-level positions (titles, paragraphs, tables, figures, form fields) with page coordinates for PDF and DOCX files. Only supported with inference=false.",
		})
		flags = append(flags, flagSchema{
			Name:        "search-results",
			Type:        "boolean",
			Required:    false,
			Location:    "body",
			Description: "When inference=true, include the raw search result chunks that were used as context for the LLM response. Defaults to false. Always true when inference=false.",
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
			Description: "Successful Response  -  returns JSON when `stream: false`, or SSE event stream when `stream: true`.",
		})

		schema := map[string]any{
			"command":     "collection-v2",
			"description": "Query Collection",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/query",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", queryCollectionV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/query", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", queryCollectionV2Flags.xOrganizationId)
	}
	if cmd.Flags().Changed("idempotency-key") {
		req.Headers["Idempotency-Key"] = fmt.Sprintf("%v", queryCollectionV2Flags.idempotencyKey)
	}

	// Request body
	bodyMap := map[string]any{}
	if queryCollectionV2Flags.body != "" {
		if err := json.Unmarshal([]byte(queryCollectionV2Flags.body), &bodyMap); err != nil {
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
	if cmd.Flags().Changed("query") {
		bodyMap["query"] = queryCollectionV2Flags.query
	}
	if cmd.Flags().Changed("inference") {
		bodyMap["inference"] = queryCollectionV2Flags.inference
	}
	if cmd.Flags().Changed("stream") {
		bodyMap["stream"] = queryCollectionV2Flags.stream
	}
	if cmd.Flags().Changed("top-k") {
		bodyMap["top_k"] = queryCollectionV2Flags.topK
	}
	if cmd.Flags().Changed("rerank") {
		bodyMap["rerank"] = queryCollectionV2Flags.rerank
	}
	if cmd.Flags().Changed("custom-prompt") {
		bodyMap["custom_prompt"] = queryCollectionV2Flags.customPrompt
	}
	if cmd.Flags().Changed("include-bbox") {
		bodyMap["include_bbox"] = queryCollectionV2Flags.includeBbox
	}
	if cmd.Flags().Changed("search-results") {
		bodyMap["search_results"] = queryCollectionV2Flags.searchResults
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
