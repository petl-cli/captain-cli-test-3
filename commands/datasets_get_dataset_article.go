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

var datasetsGetDatasetArticleCmd = &cobra.Command{
	Use:   "get-dataset-article",
	Short: "Get Dataset Article",
	RunE:  runDatasetsGetDatasetArticle,
}

var datasetsGetDatasetArticleFlags struct {
	dataset string
	url     string
}

func init() {
	datasetsGetDatasetArticleCmd.Flags().StringVar(&datasetsGetDatasetArticleFlags.dataset, "dataset", "", "The dataset to get articles from. Contact your Account Executive for available datasets.")
	datasetsGetDatasetArticleCmd.MarkFlagRequired("dataset")
	datasetsGetDatasetArticleCmd.Flags().StringVar(&datasetsGetDatasetArticleFlags.url, "url", "", "Full URL of the article to get, appended to the path. Must match the dataset's domain.")
	datasetsGetDatasetArticleCmd.MarkFlagRequired("url")

	datasetsCmd.AddCommand(datasetsGetDatasetArticleCmd)
}

func runDatasetsGetDatasetArticle(cmd *cobra.Command, args []string) error {
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
			Name:        "dataset",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "The dataset to get articles from. Contact your Account Executive for available datasets.",
		})
		flags = append(flags, flagSchema{
			Name:        "url",
			Type:        "string",
			Required:    true,
			Location:    "path",
			Description: "Full URL of the article to get, appended to the path. Must match the dataset's domain.",
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
			Description: "Article retrieved successfully",
		})
		responses = append(responses, responseSchema{
			Status:      "400",
			ContentType: "",
			Description: "Invalid request - unknown dataset or URL domain mismatch",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Unauthorized - invalid or missing API key",
		})
		responses = append(responses, responseSchema{
			Status:      "403",
			ContentType: "",
			Description: "Forbidden. Possible reasons: - API key does not belong to the specified organization - Trial expired and no payment method on file (`error: payment_required`) - Trial credits exhausted and no payment method on file (`error: credits_exhausted`)  For paid plans (Starter, Growth, Enterprise), requests are never blocked for credit usage  -  overage is billed at the end of the billing period.",
		})
		responses = append(responses, responseSchema{
			Status:      "503",
			ContentType: "",
			Description: "Dataset service unavailable",
		})

		schema := map[string]any{
			"command":     "get-dataset-article",
			"description": "Get Dataset Article",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/{dataset}/{url}",
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
	pathParams["dataset"] = fmt.Sprintf("%v", datasetsGetDatasetArticleFlags.dataset)
	pathParams["url"] = fmt.Sprintf("%v", datasetsGetDatasetArticleFlags.url)

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/{dataset}/{url}", pathParams),
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
