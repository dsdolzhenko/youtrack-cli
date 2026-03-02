package format

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

const minSeparatorWidth = 53

func separator(width int) string {
	if width < minSeparatorWidth {
		width = minSeparatorWidth
	}
	return strings.Repeat("─", width)
}

func formatMillis(ms int64) string {
	return time.UnixMilli(ms).UTC().Format("2006-01-02 15:04")
}

func terminalWidth() int {
	return 80
}

func fieldWidth(fields []youtrack.CustomField) int {
	w := 8
	for _, f := range fields {
		if len(f.Name) > w {
			w = len(f.Name)
		}
	}
	return w
}

func Issue(w io.Writer, issue *youtrack.Issue) {
	header := issue.ID + "  " + issue.Summary
	sepWidth := terminalWidth()
	if len(header) > sepWidth {
		sepWidth = len(header)
	}

	fmt.Fprintf(w, "%s\n", header)
	fmt.Fprintf(w, "%s\n", separator(sepWidth))

	fmt.Fprintf(w, "%-9s: %s\n", "Reporter", issue.Reporter.Login)
	fmt.Fprintf(w, "%-9s: %s\n", "Created", formatMillis(issue.Created))
	fmt.Fprintf(w, "%-9s: %s\n", "Updated", formatMillis(issue.Updated))
	if issue.Resolved != 0 {
		fmt.Fprintf(w, "%-9s: %s\n", "Resolved", formatMillis(issue.Resolved))
	}

	if len(issue.CustomFields) > 0 {
		fmt.Fprintf(w, "\n")
		fw := fieldWidth(issue.CustomFields)
		for _, cf := range issue.CustomFields {
			fmt.Fprintf(w, "%-*s : %s\n", fw, cf.Name, cf.Value)
		}
	}

	if issue.Description != "" {
		fmt.Fprintf(w, "\nDescription:\n")
		for _, line := range strings.Split(issue.Description, "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
	}
}

func IssueComments(w io.Writer, comments []youtrack.Comment) {
	if len(comments) == 0 {
		fmt.Fprintf(w, "\nNo comments.\n")
		return
	}
	fmt.Fprintf(w, "\nComments (%d):\n", len(comments))
	sep := separator(terminalWidth())
	for _, c := range comments {
		fmt.Fprintf(w, "%s\n", sep)
		fmt.Fprintf(w, "%s  %s\n", c.Author.Login, formatMillis(c.Created))
		fmt.Fprintf(w, "\n")
		for _, line := range strings.Split(c.Text, "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
	}
	fmt.Fprintf(w, "%s\n", sep)
}

func customFieldValue(issue youtrack.Issue, name string) string {
	for _, cf := range issue.CustomFields {
		if cf.Name == name {
			return cf.Value
		}
	}
	return ""
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func IssueList(w io.Writer, issues []youtrack.Issue) {
	const summaryMax = 50

	idW, stateW := len("ID"), len("STATE")
	for _, issue := range issues {
		if n := len(issue.ID); n > idW {
			idW = n
		}
		if s := customFieldValue(issue, "State"); len(s) > stateW {
			stateW = len(s)
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n", idW, "ID", summaryMax, "SUMMARY", stateW, "STATE", "ASSIGNEE")
	for _, issue := range issues {
		state := customFieldValue(issue, "State")
		assignee := customFieldValue(issue, "Assignee")
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n",
			idW, issue.ID,
			summaryMax, truncate(issue.Summary, summaryMax),
			stateW, state,
			assignee,
		)
	}
}
