# enexai

CLI tool that automates performance review evidence writing. Reads a CSV of
competencies, uses a self-hosted LLM to classify which skills match a piece of
feedback, analyze writing style, summarize evidence, and detect coverage gaps.

## Quick Start

```bash
cp config.ini.example config.ini   # edit with your Ollama endpoint
go build ./...
go run . --config config.ini --csv skills.csv
```

The tool prompts for a feedback title and evidence text, then classifies it
against every competency in the CSV:

```
Feedback title: API Gateway Migration (2024)
Evidence: Migrated the legacy REST gateway to a new architecture using OpenAPI specs. Evaluated three routing libraries, benchmarked latency, and implemented the chosen solution with automated integration tests.

3 skill(s) matched:
  - [fit] Evaluate and select ICT solutions
  - [fit] Implement and test software components
  - [weak_fit] Plan ICT project tasks according to process model
```

Use `--debug` for per-skill classification details.

## Self-Hosted Models

enexai is designed for self-hosted models via [Ollama](https://ollama.com).
We tested several models and recommend:

| Model | Params | Role | Notes |
|-------|--------|------|-------|
| **granite3.3:8b** | 8B | Classifier | Good instruction following, ~1s/call. Best balance of speed and accuracy. |
| **qwen3:8b** | 8B | Writer | Best quality output. Follows prompt structure precisely. Slower (~5-40s/call). |
| granite3.3:2b | 2B | - | Fast (~0.5s) but too conservative. Misses obvious matches. |
| mistral:7b | 7B | - | Tends to hallucinate technologies not mentioned in evidence. |
| llama3.2 | 3B | - | Poor instruction following. Ignores classification rules. |

Configure models in `config.ini`:

```ini
[ai]
api_url = https://your-ollama-host/v1
api_key = sk-anything

[models]
classifier = granite3.3:8b
writer = qwen3:8b
```

## Build

```bash
go build ./...          # compile
go vet ./...            # static analysis
go test ./... -v        # run tests
```

## License

Private.
