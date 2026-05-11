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

var indexingIndexYoutubeV2Cmd = &cobra.Command{
	Use:   "index-youtube-v2",
	Short: "Index YouTube Transcript",
	RunE:  runIndexingIndexYoutubeV2,
}

var indexingIndexYoutubeV2Flags struct {
	xOrganizationId string
	collectionName  string
	idempotencyKey  string
	url             string
	urls            []string
	languages       []string
	body            string
}

func init() {
	indexingIndexYoutubeV2Cmd.Flags().StringVar(&indexingIndexYoutubeV2Flags.xOrganizationId, "x-organization-id", "", "The organization ID to scope the request")
	indexingIndexYoutubeV2Cmd.MarkFlagRequired("x-organization-id")
	indexingIndexYoutubeV2Cmd.Flags().StringVar(&indexingIndexYoutubeV2Flags.collectionName, "collection-name", "", "Name of the collection to index into")
	indexingIndexYoutubeV2Cmd.MarkFlagRequired("collection-name")
	indexingIndexYoutubeV2Cmd.Flags().StringVar(&indexingIndexYoutubeV2Flags.idempotencyKey, "idempotency-key", "", "UUID for request deduplication")
	indexingIndexYoutubeV2Cmd.Flags().StringVar(&indexingIndexYoutubeV2Flags.url, "url", "", "A single YouTube video URL. Supported formats: youtube.com/watch?v=, youtu.be/, youtube.com/shorts/. Provide either 'url' or 'urls', not both.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexYoutubeV2Cmd.Flags().StringSliceVar(&indexingIndexYoutubeV2Flags.urls, "urls", nil, "An array of YouTube video URLs to index (max 20). Provide either 'url' or 'urls', not both.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexYoutubeV2Cmd.Flags().StringSliceVar(&indexingIndexYoutubeV2Flags.languages, "languages", nil, "Preferred transcript languages in priority order (ISO 639-1 codes). Defaults to English. Only specify if you need a non-English transcript (e.g., ['fr', 'de']). Falls back to auto-generated captions if manual transcript unavailable.")
	// Note: body fields are not MarkFlagRequired — --body JSON satisfies them too.
	indexingIndexYoutubeV2Cmd.Flags().StringVar(&indexingIndexYoutubeV2Flags.body, "body", "", "Full request body as JSON. Individual body flags override matching keys in this JSON.")

	indexingCmd.AddCommand(indexingIndexYoutubeV2Cmd)
}

func runIndexingIndexYoutubeV2(cmd *cobra.Command, args []string) error {
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
			Description: "A single YouTube video URL. Supported formats: youtube.com/watch?v=, youtu.be/, youtube.com/shorts/. Provide either 'url' or 'urls', not both.",
		})
		flags = append(flags, flagSchema{
			Name:        "urls",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "An array of YouTube video URLs to index (max 20). Provide either 'url' or 'urls', not both.",
		})
		flags = append(flags, flagSchema{
			Name:        "languages",
			Type:        "array",
			Required:    false,
			Location:    "body",
			Description: "Preferred transcript languages in priority order (ISO 639-1 codes). Defaults to English. Only specify if you need a non-English transcript (e.g., ['fr', 'de']). Falls back to auto-generated captions if manual transcript unavailable.",
		})
		flags = append(flags, flagSchema{
			Name:        "custom-metadata",
			Type:        "object",
			Required:    false,
			Location:    "body",
			Description: "Custom metadata to attach to all indexed chunks. Keys must be strings. Values: str, int, float, bool, or array of strings.",
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
			"command":     "index-youtube-v2",
			"description": "Index YouTube Transcript",
			"http": map[string]any{
				"method": "POST",
				"path":   "/v2/collections/{collection_name}/index/youtube",
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
	pathParams["collection_name"] = fmt.Sprintf("%v", indexingIndexYoutubeV2Flags.collectionName)

	req := &httpclient.Request{
		Method:      "POST",
		Path:        httpclient.SubstitutePath("/v2/collections/{collection_name}/index/youtube", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters

	// Header parameters
	if cmd.Flags().Changed("x-organization-id") {
		req.Headers["X-Organization-ID"] = fmt.Sprintf("%v", indexingIndexYoutubeV2Flags.xOrganizationId)
	}
	if cmd.Flags().Changed("idempotency-key") {
		req.Headers["Idempotency-Key"] = fmt.Sprintf("%v", indexingIndexYoutubeV2Flags.idempotencyKey)
	}

	// Request body
	bodyMap := map[string]any{}
	if indexingIndexYoutubeV2Flags.body != "" {
		if err := json.Unmarshal([]byte(indexingIndexYoutubeV2Flags.body), &bodyMap); err != nil {
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
		bodyMap["url"] = indexingIndexYoutubeV2Flags.url
	}
	if cmd.Flags().Changed("urls") {
		bodyMap["urls"] = indexingIndexYoutubeV2Flags.urls
	}
	if cmd.Flags().Changed("languages") {
		bodyMap["languages"] = indexingIndexYoutubeV2Flags.languages
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
