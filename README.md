# YouTrack CLI (yt)

CLI utility to access YouTrack.

## Installation

Using homebrew:

```sh
brew install dsdolzhenko/tools/youtrack-cli
```

With `go install`:

```sh
go install github.com/dsdolzhenko/youtrack-cli@latest
```

Or build from source:

```sh
go build -o yt .
```

## Configuration

Configuration is resolved in priority order: **flag > environment variable > config file**.

| Setting    | Flag      | Environment variable | Config file key |
|------------|-----------|----------------------|-----------------|
| Server URL | `--url`   | `YOUTRACK_URL`       | `url`           |
| API token  | `--token` | `YOUTRACK_TOKEN`     | `token`         |

Config file location: `~/.config/youtrack/config.json`

```json
{
  "url": "https://your-instance.youtrack.cloud",
  "token": "perm:your-permanent-token"
}
```

## Claude Code integration

`yt` can be used by Claude Code agents in any project. Run once after installing:

```sh
yt setup
```

This installs a skill into `~/.claude/skills/youtrack/`. After that, you can ask Claude things like:

- "look up issue SP-42"
- "search YouTrack for open bugs assigned to me"
- "/youtrack project: SP State: Open"

Claude will invoke `yt` automatically and summarise the results.

## Commands

### `yt issue <ID>`

Fetch and display a single issue.

```
$ yt issue SP-42
SP-42  Fix login redirect loop
─────────────────────────────────────────────────────────────────────────────────
Reporter : jane.doe
Created  : 2024-11-01 14:32
Updated  : 2024-11-02 10:00
Resolved : 2024-11-03 09:15

Priority : Critical
State    : Fixed
Assignee : John Doe (john.doe)
Sprint   : Sprint 12

Description:
  When a user logs in and is redirected...
```

### `yt article <ID>`

Fetch and display a single knowledge base article.

```
$ yt article A-7
A-7  Deployment guide
─────────────────────────────────────────────────────────────────────────────────
Project  : MyProject (PROJ)
Reporter : jane.doe
Created  : 2024-10-15 09:00
Updated  : 2024-10-20 11:30

Content:
  ## Prerequisites
  ...
```

### `yt search <query>`

Search issues using the YouTrack query language.

```
$ yt search "project: SP State: Open" --top 20
ID        SUMMARY                                        STATE          ASSIGNEE
SP-42     Fix login redirect loop                        Open           John Doe
SP-43     Add dark mode                                  In Progress    Jane Doe
```

**Flags:**
- `--top N` — maximum number of results (default: 50)

### `yt command <ID> "<command>"`

Apply a [YouTrack command](https://www.jetbrains.com/help/youtrack/server/youtrack-command-syntax-reference.html) to an issue.

```
$ yt command SP-42 "state Fixed"
Command applied to SP-42

$ yt command SP-42 "for me"
Command applied to SP-42

$ yt command SP-42 "tag needs-review" --comment "Please review" --silent
Command applied to SP-42
```

**Flags:**
- `--comment "..."` — attach a comment alongside the command
- `--silent` — apply the command without sending YouTrack notifications
