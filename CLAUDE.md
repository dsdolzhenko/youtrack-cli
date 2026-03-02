# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```sh
go build ./...           # build
go test ./...            # run all tests
go test ./internal/...   # run a specific package
go test -run TestName ./internal/youtrack/  # run a single test
go run . --help          # inspect CLI commands
```

## Architecture

The app is a read-only CLI for YouTrack built with Cobra. Data flows strictly in one direction:

```
cmd/ → internal/youtrack/ → internal/client/ → YouTrack API
cmd/ → internal/format/   (renders youtrack types to io.Writer)
```

**`internal/client`** — thin HTTP wrapper. Sets auth headers, builds URLs, returns `*http.Response` on 2xx or an error otherwise. Callers own JSON decoding.

**`internal/youtrack`** — API types and fetch functions (`GetIssue`, `SearchIssues`, `GetArticle`, `SearchArticles`). Issues require a two-pass JSON decode: the `customFields` array is first captured as `[]json.RawMessage` via an intermediate `issueRaw` struct, then passed to `DecodeCustomFields` which dispatches on the `$type` discriminator field to produce `[]CustomField{Name, Value string}`. Articles decode in a single pass (no polymorphic fields).

**`internal/format`** — pure rendering, no network calls. List formatters (`IssueList`, `ArticleList`) compute column widths dynamically from the data before printing. All functions write to an `io.Writer`.

**`cmd/`** — Cobra wiring only. Shared run functions (`runGetIssue`, `runSearchIssues`, `runGetArticle`, `runSearchArticles`) live in `cmd/issues.go` and `cmd/articles.go` and are called by both the full subcommands (`yt issues get`, `yt articles search`, …) and the top-level shortcuts (`yt issue`, `yt article`).

## Configuration

Priority: flag > env var > config file (`~/.config/youtrack-cli/config.json`).

```json
{ "url": "https://your-instance.youtrack.cloud", "token": "perm:..." }
```

Env vars: `YOUTRACK_URL`, `YOUTRACK_TOKEN`.

## Testing

Tests use `net/http/httptest.NewServer` exclusively — no external dependencies. The `newTestClient` helper shared across `issue_test.go` and `article_test.go` lives in `issue_test.go` (same `package youtrack` test binary).
