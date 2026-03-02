# youtrack-cli (yt)

CLI utility to access YouTrack that is designed to be used both by humans and LLMs.

## Installation

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

Config file location: `~/.config/youtrack-cli/config.json`

```json
{
  "url": "https://your-instance.youtrack.cloud",
  "token": "perm:your-permanent-token"
}
```

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
