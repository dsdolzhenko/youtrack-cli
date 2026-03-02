package youtrack

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/dsdolzhenko/youtrack-cli/internal/client"
)

type Project struct {
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
}

type Article struct {
	ID        string  `json:"idReadable"`
	Summary   string  `json:"summary"`
	Content   string  `json:"content"`
	Created   int64   `json:"created"`
	Updated   int64   `json:"updated"`
	Reporter  User    `json:"reporter"`
	Project   Project `json:"project"`
}

const articleFields = "id,idReadable,summary,content,created,updated," +
	"reporter(login,fullName),project(shortName,name)"

const searchArticleFields = "id,idReadable,summary,created,updated," +
	"reporter(login,fullName),project(shortName,name)"

func GetArticle(c *client.Client, id string) (*Article, error) {
	params := url.Values{}
	params.Set("fields", articleFields)

	resp, err := c.Get("/api/articles/"+id, params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: get article %s: %w", id, err)
	}
	defer resp.Body.Close()

	var article Article
	if err := json.NewDecoder(resp.Body).Decode(&article); err != nil {
		return nil, fmt.Errorf("youtrack: decode article %s: %w", id, err)
	}

	return &article, nil
}

func SearchArticles(c *client.Client, query string, top int) ([]Article, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("$top", fmt.Sprintf("%d", top))
	params.Set("fields", searchArticleFields)

	resp, err := c.Get("/api/articles", params)
	if err != nil {
		return nil, fmt.Errorf("youtrack: search articles: %w", err)
	}
	defer resp.Body.Close()

	var articles []Article
	if err := json.NewDecoder(resp.Body).Decode(&articles); err != nil {
		return nil, fmt.Errorf("youtrack: decode article search results: %w", err)
	}

	return articles, nil
}
