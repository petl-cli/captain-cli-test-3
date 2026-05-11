package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// agentInstructionsCmd prints the llms.txt content at runtime so agents can
// include it in their system prompt without needing a separate file:
//
//	INSTRUCTIONS=$(captain-api-v2 agent-instructions)
var agentInstructionsCmd = &cobra.Command{
	Use:   "agent-instructions",
	Short: "Print machine-readable instructions for AI agents (llms.txt format)",
	Long: `Prints a complete description of this CLI's commands, flags, exit codes,
and usage patterns optimised for inclusion in an AI agent's system prompt.

Example:
  # Include in Claude Code context:
  captain-api-v2 agent-instructions > CLAUDE.md

  # Capture inline:
  INSTRUCTIONS=$(captain-api-v2 agent-instructions)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print(agentInstructionsContent)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentInstructionsCmd)
}

// agentInstructionsContent is the full llms.txt baked into the binary at build time.
// Regenerate by re-running the CLI generator against the same OpenAPI spec.
const agentInstructionsContent = `# captain-api-v2

Captain API v2 - Agentic Data Indexing & Retrieval Platform. RESTful API with improved resource-based URLs.

This file is the agent-facing overview of the ` + "`" + `captain-api-v2` + "`" + ` CLI. It explains what the tool does and how to use it well — *not* every flag of every command. For per-command details run ` + "`" + `captain-api-v2 <command> --help` + "`" + ` or ` + "`" + `captain-api-v2 <command> --schema` + "`" + ` (returns JSON).

## Install

The binary is ` + "`" + `captain-api-v2` + "`" + `. Build from source or download a release.

## Authentication

- Bearer token — set ` + "`" + `CAPTAIN_API_V2_BEARER_TOKEN` + "`" + ` or pass ` + "`" + `--bearer-token <token>` + "`" + `

Run ` + "`" + `captain-api-v2 configure` + "`" + ` to see how to set credentials.

## Commands

captain-api-v2 groups operations by resource. One or two examples per group below — for the full set of commands and flags use ` + "`" + `captain-api-v2 <group> --help` + "`" + `.

### ` + "`" + `collections` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 collections change-collection-environment-v2 --x-organization-id <string> --collection-name <string> --body '{...}'    # Change Collection Environment
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 collections create-collection-v2 --x-organization-id <string> --collection-name <string>    # Create Collection
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 collections --help` + "`" + `

### ` + "`" + `companies` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 companies active-investors --x-organization-id <string> --company-id <string>    # Get Company Active Investors
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 companies bio --x-organization-id <string> --company-id <string>    # Get Company Bio
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 companies --help` + "`" + `

### ` + "`" + `credit-analysis` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 credit-analysis agreement-detail --x-organization-id <string> --company <string>    # Get Agreement Detail
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 credit-analysis agreements-search --x-organization-id <string>    # Search Credit Agreements
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 credit-analysis --help` + "`" + `

### ` + "`" + `datasets` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 datasets batch-search --x-organization-id <string> --body '{...}'    # Batch Search Articles
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 datasets get-dataset-article --x-organization-id <string> --dataset <string> --url <string>    # Get Dataset Article
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 datasets --help` + "`" + `

### ` + "`" + `deals` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 deals bio --x-organization-id <string> --id <string>    # Get Deal Bio
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 deals cap-table --x-organization-id <string> --id <string>    # Get Deal Cap Table
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 deals --help` + "`" + `

### ` + "`" + `funds` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 funds active-investments --x-organization-id <string> --fund-id <string>    # Get Fund Active Investments
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 funds benchmark --x-organization-id <string> --fund-id <string>    # Get Fund Benchmark Comparison
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 funds --help` + "`" + `

### ` + "`" + `general` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 general entity-affiliates --x-organization-id <string> --entity-id <string>    # Get Entity Affiliates
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 general entity-locations --x-organization-id <string> --entity-id <string>    # Get Entity Locations
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 general --help` + "`" + `

### ` + "`" + `indexing` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 indexing index-azure-container-v2 --x-organization-id <string> --collection-name <string> --body '{...}'    # Index Azure Container
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 indexing index-azure-directory-v2 --x-organization-id <string> --collection-name <string> --body '{...}'    # Index Azure Directory
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 indexing --help` + "`" + `

### ` + "`" + `investors` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 investors active-investments --x-organization-id <string> --id <string>    # Get Investor Active Investments
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 investors bio --x-organization-id <string> --id <string>    # Get Investor Bio
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 investors --help` + "`" + `

### ` + "`" + `jobs` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 jobs delete-job-v2 --x-organization-id <string> --job-id <string>    # Delete Job
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 jobs get-job-status-v2 --x-organization-id <string> --job-id <string>    # Get Job Status
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 jobs --help` + "`" + `

### ` + "`" + `limited-partners` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 limited-partners lps-allocations-actual --x-organization-id <string> --lp-id <string>    # Get LP Actual Allocations
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 limited-partners lps-allocations-target --x-organization-id <string> --lp-id <string>    # Get LP Target Allocations
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 limited-partners --help` + "`" + `

### ` + "`" + `patents` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 patents get-by-id --x-organization-id <string> --id <string>    # Get Patent Details
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 patents get-file --x-organization-id <string> --entity-id <string>    # Get Patent File
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 patents --help` + "`" + `

### ` + "`" + `people` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 people bio --x-organization-id <string> --person-id <string>    # Enrich Person
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 people search --x-organization-id <string> --q <string>    # Search People
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 people --help` + "`" + `

### ` + "`" + `query` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 query collection-v2 --x-organization-id <string> --collection-name <string> --body '{...}'    # Query Collection
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 query --help` + "`" + `

### ` + "`" + `sandbox-data` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 sandbox-data fundamentals-lookup-table-values --x-organization-id <string> --table-id <string>    # Get Lookup Table Values
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 sandbox-data fundamentals-lookup-tables --x-organization-id <string>    # Get Lookup Tables
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 sandbox-data --help` + "`" + `

### ` + "`" + `service-providers` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 service-providers bio --x-organization-id <string> --provider-id <string>    # Get Service Provider Bio
` + "`" + `` + "`" + `` + "`" + `

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 service-providers companies --x-organization-id <string> --provider-id <string>    # Get Service Provider Company Clients
` + "`" + `` + "`" + `` + "`" + `

List all commands in this group: ` + "`" + `captain-api-v2 service-providers --help` + "`" + `


## Output and parsing

All output is JSON by default and goes to stdout. Errors go to stderr as a single-line JSON object — stdout stays clean for piping.

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 <cmd> --jq <path>          # extract fields without jq installed (GJSON syntax)
captain-api-v2 <cmd> -o yaml              # change format: json (default), yaml, table, compact, raw, pretty
captain-api-v2 <cmd> --dry-run            # print the HTTP request without sending it
captain-api-v2 <cmd> --schema             # JSON schema of inputs and outputs for that command
` + "`" + `` + "`" + `` + "`" + `

GJSON path examples:
- ` + "`" + `--jq id` + "`" + ` — scalar
- ` + "`" + `--jq items.#.id` + "`" + ` — every id from an array
- ` + "`" + `--jq "items.#(active==true)#"` + "`" + ` — filter array by condition

## Exit codes

Branch on ` + "`" + `$?` + "`" + ` rather than parsing stderr:

| Code | Meaning |
|------|---------|
| 0 | success |
| 1 | unknown error |
| 2 | auth failed (401 / 403) |
| 3 | not found (404) |
| 4 | validation error (400 / 422) |
| 5 | rate limited (429) |
| 6 | server error (5xx) |
| 7 | network error |

Error JSON shape on stderr:
` + "`" + `` + "`" + `` + "`" + `json
{"error":true,"code":"not_found","status":404,"message":"...","exit_code":3}
` + "`" + `` + "`" + `` + "`" + `

## Common workflows

### List then fetch one
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(captain-api-v2 <group> list --jq "items.0.id")
captain-api-v2 <group> get --id "$ID"
` + "`" + `` + "`" + `` + "`" + `

### Capture a created ID
` + "`" + `` + "`" + `` + "`" + `bash
ID=$(captain-api-v2 <group> create [flags] --jq id)
` + "`" + `` + "`" + `` + "`" + `

### Safe destructive call
` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 <group> delete --id X --dry-run    # inspect first
captain-api-v2 <group> delete --id X              # exit 0 = deleted, 3 = not found
` + "`" + `` + "`" + `` + "`" + `

### Branch on exit code
` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 <cmd> [flags]
case $? in
  0) : success ;;
  2) echo "fix credentials" ;;
  3) echo "not found, skip" ;;
  5) sleep 10 ; retry ;;
esac
` + "`" + `` + "`" + `` + "`" + `

## Discovering more

` + "`" + `` + "`" + `` + "`" + `bash
captain-api-v2 --help                              # top-level groups
captain-api-v2 <group> --help                      # commands in a group
captain-api-v2 <group> <cmd> --help                # flags + description
captain-api-v2 <group> <cmd> --schema              # JSON schema
captain-api-v2 agent-instructions                  # this file, embedded in the binary
` + "`" + `` + "`" + `` + "`" + `

`
