# YouTrack Query Language Reference

## Common filter keywords

| Keyword      | Example                        | Notes                              |
|--------------|--------------------------------|------------------------------------|
| `project:`   | `project: SP`                  | Project short name                 |
| `State:`     | `State: Open`                  | Issue state name                   |
| `Priority:`  | `Priority: Critical`           | Priority name                      |
| `assignee:`  | `assignee: me`                 | Username or `me` for current user  |
| `reporter:`  | `reporter: john.doe`           | Username                           |
| `Type:`      | `Type: Bug`                    | Issue type                         |
| `Sprint:`    | `Sprint: "Sprint 5"`           | Sprint name (quote if multi-word)  |
| `tag:`       | `tag: backend`                 | Tag name                           |

## Date ranges

```
created: 2024-01-01 .. today
updated: -1w .. today
resolved: -30d .. today
```

Date shortcuts: `today`, `-1d`, `-1w`, `-1m`, `-1y`

## Negation

Prefix a value with `-` to exclude it:

```
State: -Resolved
assignee: -me
Priority: -Minor
```

## Free-text search

Any term without a keyword prefix searches in summary and description:

```
login error 500
```

Combine with filters:

```
login error project: SP State: Open
```

## Combining filters

Space-separated terms are ANDed:

```
project: SP assignee: me State: Open Priority: Critical
```

## Sort order

`yt` uses YouTrack's default sort order (updated descending). There is no CLI flag to change sort order.

## Examples

```bash
# Open bugs in project SP assigned to me
yt issues search "project: SP Type: Bug State: Open assignee: me"

# All critical issues updated in the last week
yt issues search "Priority: Critical updated: -1w .. today State: -Resolved"

# Search for login-related issues
yt issues search "login State: Open"

# Issues in a specific sprint
yt issues search "project: SP Sprint: \"Sprint 10\" State: -Resolved"
```
