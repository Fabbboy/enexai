# AGENTS.md — enexai

CLI tool that automates performance review evidence writing. Reads a CSV of
competencies, uses a self-hosted LLM (Ollama) to classify skill fit, analyze
writing style, summarize evidence, and detect coverage gaps.

## Build & Run

```bash
go build ./...          # compile
go vet ./...            # static analysis
go test ./... -v        # all tests
go test ./impl -run TestParse  # single test by name
go run . --config config.ini --csv skills.csv  # run (needs real config + CSV)
```

AI calls hit a self-hosted GPU endpoint (~5-30s per call). Use generous
timeouts when running end-to-end.

## Project Layout

```
main.go              # Cobra CLI entry — DO NOT modify unless adding flags
impl/
  ai.go              # aiClient struct + thin helpers (userMsg, systemMsg, parse, etc.)
  skills.go          # Skill type, all AI agent functions, JSON schemas, //go:embed
  config.go          # INI config loader (AiConfig, ModelConfig)
  driver.go          # Orchestration prototype — hardcoded test data is intentional
  prompts/           # Embedded .txt prompt files
    fits_skill.txt
    style_analysis.txt
    evidence_analysis.txt
    coverage_detection.txt
config.ini.example   # Template for config — real config.ini is gitignored
skills.csv           # Real competency data — gitignored
```

Everything lives in package `impl`. Do not restructure packages.

## Architecture Rules

### AI Helpers (ai.go)

`aiClient` is a thin struct holding ctx/client/model/logger. Its `Send()`
method calls `client.Responses.New()` directly and adds timing/token logging.

- `client.Responses.New()` MUST remain the visible API call inside Send
- Do NOT add abstractions that wrap or hide the API call further
- Helpers exist to reduce openai-go struct boilerplate, not to abstract:
  `userMsg`, `systemMsg`, `inputItems`, `jsonSchemaFormat`, `parse[T]`

### AI Functions (skills.go)

Each function builds `ResponseNewParams` directly and calls `client.Send()`.
Pattern:

```go
func SomeAgent(client aiClient, skill *Skill) (*Result, error) {
    instructions := somePrompt + "\n\n" + skill.FormatContext()
    params := responses.ResponseNewParams{
        Instructions: openai.String(instructions),
        Input:        inputItems(userMsg(...)),
        Text:         jsonSchemaFormat("result_name", resultSchema),
    }
    resp, err := client.Send(params)
    // ... parse[ResultType](resp)
}
```

- System context goes in `Instructions`, user content goes in `Input`
- Every result type has a companion `var ...Schema = map[string]any{...}`
  defined right above its function
- Use `parse[T]` for JSON unmarshaling — do not call json.Unmarshal directly

### Prompts (impl/prompts/)

All prompt files follow the same structure:

```
ROLE
<one line>

TASK
<what to do>

OUTPUT
<JSON shape description>

CONSTRAINTS
<bullet list>
```

Keep this structure when adding new prompts. Do not deviate.

## Code Style

### General

- Go 1.24, `gofmt` formatting, no external linter configured
- Prefer simple, pragmatic code — no unnecessary abstractions
- No dead code — remove what you don't use
- No `as any`-style suppressions or ignoring errors

### Naming

- Exported: `FitsSkill`, `LoadConfig`, `Skill`, `FitsResult`
- Unexported internals: `aiClient`, `userMsg`, `parse`, `fitsResultSchema`
- Typed string enums: `type Fitness string` with const block
- JSON tags on result structs, CSV tags on `Skill`, INI tags on config

### Imports

Group in this order (separated by blank lines):

1. Standard library
2. External dependencies

```go
import (
    _ "embed"
    "os"
    "strings"

    "github.com/gocarina/gocsv"
    "github.com/openai/openai-go/v3"
)
```

### Error Handling

- Return errors directly: `return nil, err`
- `Send()` logs errors before returning — callers should not double-log
- Do not swallow errors with empty catch blocks

### Testing

- Standard `testing` package only — no testify, gomock, etc.
- Do NOT call the real AI API in tests
- Use `t.TempDir()` for temp files, inline data for test fixtures
- Table-driven tests where it makes sense

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/openai/openai-go/v3` | OpenAI Responses API client |
| `github.com/gocarina/gocsv` | CSV parsing into structs |
| `gopkg.in/ini.v1` | INI config file parsing |
| `github.com/spf13/cobra` | CLI framework |
| `github.com/lmittmann/tint` | Colored structured logging |

## Sensitive Files

These are gitignored and contain real data — do not commit:

- `config.ini` — API URL, keys, model names
- `skills.csv` / `*.xlsx` — real competency portfolio data

Use `config.ini.example` as reference for config structure.

## What NOT to Do

- Do NOT modify `main.go` beyond adding new CLI flags
- Do NOT change the driver flow/logic (prototype with hardcoded test data)
- Do NOT add an evidence writer (step 3) — that's future work
- Do NOT restructure packages or rename `impl`
- Do NOT add unnecessary abstractions — thin helpers only
- Do NOT add comments that restate what the code does
