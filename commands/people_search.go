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

var peopleSearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search People",
	RunE:  runPeopleSearch,
}

var peopleSearchFlags struct {
	q        string
	company  string
	title    string
	location string
	limit    int
	offset   int
}

func init() {
	peopleSearchCmd.Flags().StringVar(&peopleSearchFlags.q, "q", "", "Person name or search query")
	peopleSearchCmd.MarkFlagRequired("q")
	peopleSearchCmd.Flags().StringVar(&peopleSearchFlags.company, "company", "", "Filter by current company")
	peopleSearchCmd.Flags().StringVar(&peopleSearchFlags.title, "title", "", "Filter by job title")
	peopleSearchCmd.Flags().StringVar(&peopleSearchFlags.location, "location", "", "Filter by location")
	peopleSearchCmd.Flags().IntVar(&peopleSearchFlags.limit, "limit", 0, "Maximum results per page (1-500). For limits above 100, query expansion is used automatically.")
	peopleSearchCmd.Flags().IntVar(&peopleSearchFlags.offset, "offset", 0, "Number of results to skip for pagination. Use with limit to page through results.")

	peopleCmd.AddCommand(peopleSearchCmd)
}

func runPeopleSearch(cmd *cobra.Command, args []string) error {
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
			Name:        "q",
			Type:        "string",
			Required:    true,
			Location:    "query",
			Description: "Person name or search query",
		})
		flags = append(flags, flagSchema{
			Name:        "company",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by current company",
		})
		flags = append(flags, flagSchema{
			Name:        "title",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by job title",
		})
		flags = append(flags, flagSchema{
			Name:        "location",
			Type:        "string",
			Required:    false,
			Location:    "query",
			Description: "Filter by location",
		})
		flags = append(flags, flagSchema{
			Name:        "limit",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Maximum results per page (1-500). For limits above 100, query expansion is used automatically.",
		})
		flags = append(flags, flagSchema{
			Name:        "offset",
			Type:        "integer",
			Required:    false,
			Location:    "query",
			Description: "Number of results to skip for pagination. Use with limit to page through results.",
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
			Description: "Successful response",
		})
		responses = append(responses, responseSchema{
			Status:      "401",
			ContentType: "",
			Description: "Unauthorized - Invalid or missing API key",
		})

		schema := map[string]any{
			"command":     "search",
			"description": "Search People",
			"http": map[string]any{
				"method": "GET",
				"path":   "/v2/datasets/odyssey/people/search",
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

	req := &httpclient.Request{
		Method:      "GET",
		Path:        httpclient.SubstitutePath("/v2/datasets/odyssey/people/search", pathParams),
		QueryParams: map[string]string{},
		ArrayParams: map[string][]string{},
		Headers:     map[string]string{},
	}

	// Query parameters
	if cmd.Flags().Changed("q") {
		req.QueryParams["q"] = fmt.Sprintf("%v", peopleSearchFlags.q)
	}
	if cmd.Flags().Changed("company") {
		req.QueryParams["company"] = fmt.Sprintf("%v", peopleSearchFlags.company)
	}
	if cmd.Flags().Changed("title") {
		req.QueryParams["title"] = fmt.Sprintf("%v", peopleSearchFlags.title)
	}
	if cmd.Flags().Changed("location") {
		req.QueryParams["location"] = fmt.Sprintf("%v", peopleSearchFlags.location)
	}
	if cmd.Flags().Changed("limit") {
		req.QueryParams["limit"] = fmt.Sprintf("%v", peopleSearchFlags.limit)
	}
	if cmd.Flags().Changed("offset") {
		req.QueryParams["offset"] = fmt.Sprintf("%v", peopleSearchFlags.offset)
	}

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
