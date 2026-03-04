package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
	"github.com/dsdolzhenko/youtrack-cli/internal/format"
	"github.com/dsdolzhenko/youtrack-cli/internal/youtrack"
)

var articlesCmd = &cobra.Command{
	Use:   "articles",
	Short: "Work with articles",
}

var articlesGetCmd = &cobra.Command{
	Use:   "get <ID>",
	Short: "Show a single article",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGetArticle(args[0])
	},
}

var articlesSearchTop int
var articlesSearchIDs bool

var articlesSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search articles using YouTrack query language",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSearchArticles(args[0], articlesSearchTop)
	},
}

func runGetArticle(id string) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	article, err := youtrack.GetArticle(c, id)
	if err != nil {
		return err
	}
	if jsonOutput {
		return writeJSON(article)
	}
	format.Article(os.Stdout, article)
	return nil
}

func runSearchArticles(query string, top int) error {
	if err := checkConfig(); err != nil {
		return err
	}
	c := client.New(client.Config{BaseURL: serverURL, Token: token})
	articles, err := youtrack.SearchArticles(c, query, top)
	if err != nil {
		return err
	}
	if articlesSearchIDs {
		for _, article := range articles {
			fmt.Fprintln(os.Stdout, article.ID)
		}
		return nil
	}
	if jsonOutput {
		return writeJSON(articles)
	}
	format.ArticleList(os.Stdout, articles)
	return nil
}

func init() {
	articlesSearchCmd.Flags().IntVar(&articlesSearchTop, "top", 50, "Maximum number of results")
	articlesSearchCmd.Flags().BoolVar(&articlesSearchIDs, "ids", false, "Print only article IDs, one per line")
	articlesCmd.AddCommand(articlesGetCmd)
	articlesCmd.AddCommand(articlesSearchCmd)
	rootCmd.AddCommand(articlesCmd)
}
