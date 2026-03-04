---
description: Use when the user asks about YouTrack issues, tickets, tasks, articles, or wants to search, look up, or summarise work items tracked in YouTrack. Triggers on phrases like "look up issue", "find ticket", "search YouTrack", "get article", "show me open bugs".
argument-hint: "[issue ID | article ID | search query]"
allowed-tools: Bash
user-invocable: true
---

# YouTrack skill

You have access to a `yt` CLI for YouTrack access.

## Commands

```bash
yt issue <ID>                                      # Get a single issue (e.g. yt issue SP-42)
yt article <ID>                                    # Get a single article
yt issues search "<query>" [--top N]               # Search issues (default --top 50)
yt articles search "<query>" [--top N]             # Search articles
yt command <ID> "<command>" [--comment ""] [--silent]  # Apply a YouTrack command to an issue
```

## Configuration

Priority: flag > env var > config file.

Config file: `~/.config/youtrack/config.json`
```json
{ "url": "https://your-instance.youtrack.cloud", "token": "perm:..." }
```

Env vars: `YOUTRACK_URL`, `YOUTRACK_TOKEN`

If a command fails with "server URL is required" or "token is required", tell the user to set up credentials.

## Output formats

**Single issue** — header fields (Reporter, dates), custom fields (Priority, State, Assignee, Sprint, etc.), then full description.

**Issue list** — columns: ID · SUMMARY (≤50 chars) · STATE · ASSIGNEE.

**Single article** — header (Project, Reporter, dates), then full Markdown content.

**Article list** — columns: ID · SUMMARY · PROJECT · REPORTER.

## Commands guide

Use `yt command` to apply YouTrack commands (change state, assign, tag, etc.):

```bash
yt command SP-42 "state Fixed"
yt command SP-42 "for me"
yt command SP-42 "tag needs-review" --comment "Please review" --silent
```

`--silent` suppresses YouTrack notifications. `--comment` attaches a comment alongside the command.

## Behaviour when invoked as `/youtrack $ARGUMENTS`

- `$ARGUMENTS` matches `[A-Z]+-\d+` (e.g. `SP-42`) → `yt issue $ARGUMENTS`
- Otherwise → treat as a search query → `yt issues search "$ARGUMENTS"`
- No arguments → infer intent from conversation context.

When searching issues, prefer narrow queries — add `State: -Resolved` unless the user explicitly wants resolved items.

See `references/query-language.md` for YouTrack query syntax.
