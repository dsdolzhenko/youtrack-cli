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

## Usage

```sh
yt issue SP-42                                          # view issue
yt search "project: SP State: Open"                     # search issues
yt issues attachments SP-42 --dir ./files               # download all attachments
yt issues attachments SP-42 schema.yaml                 # download one attachment
yt command SP-42 "state Fixed" --comment "done" --silent
```

Run `yt --help` or `yt <command> --help` for the full reference.

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
