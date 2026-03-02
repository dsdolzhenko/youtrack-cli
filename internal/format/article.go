package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

func Article(w io.Writer, article *youtrack.Article) {
	header := article.ID + "  " + article.Summary
	sepWidth := terminalWidth()
	if len(header) > sepWidth {
		sepWidth = len(header)
	}

	fmt.Fprintf(w, "%s\n", header)
	fmt.Fprintf(w, "%s\n", separator(sepWidth))

	fmt.Fprintf(w, "%-9s: %s (%s)\n", "Project", article.Project.Name, article.Project.ShortName)
	fmt.Fprintf(w, "%-9s: %s\n", "Reporter", article.Reporter.Login)
	fmt.Fprintf(w, "%-9s: %s\n", "Created", formatMillis(article.Created))
	fmt.Fprintf(w, "%-9s: %s\n", "Updated", formatMillis(article.Updated))

	if article.Content != "" {
		fmt.Fprintf(w, "\nContent:\n")
		for _, line := range strings.Split(article.Content, "\n") {
			fmt.Fprintf(w, "  %s\n", line)
		}
	}
}

func ArticleList(w io.Writer, articles []youtrack.Article) {
	const summaryMax = 50

	idW, projW := len("ID"), len("PROJECT")
	for _, a := range articles {
		if n := len(a.ID); n > idW {
			idW = n
		}
		if n := len(a.Project.ShortName); n > projW {
			projW = n
		}
	}

	fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n", idW, "ID", summaryMax, "SUMMARY", projW, "PROJECT", "REPORTER")
	for _, a := range articles {
		fmt.Fprintf(w, "%-*s  %-*s  %-*s  %s\n",
			idW, a.ID,
			summaryMax, truncate(a.Summary, summaryMax),
			projW, a.Project.ShortName,
			a.Reporter.Login,
		)
	}
}
