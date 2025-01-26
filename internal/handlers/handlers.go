package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
)

// Index はRSSフィードを取得して、HTMLまたはJSONで出力するハンドラーです。
func Index(c echo.Context) error {
	log.Println("Indexハンドラーが呼び出されました") // デバッグ用ログ

	// RSSフィードURL
	rssURL := "https://feed.infoq.com/"
	log.Println("RSSフィードURL:", rssURL) // デバッグ用ログ

	// パーサーの作成
	fp := gofeed.NewParser()

	// RSSフィードを取得して解析
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		log.Println("RSSフィード取得エラー:", err) // デバッグ用ログ
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("RSSフィードの取得に失敗しました: %v", err),
		})
	}

	log.Println("RSSフィード取得成功:", feed.Title) // デバッグ用ログ

	// フィード情報を整形
	feedItems := []map[string]interface{}{}
	for _, item := range feed.Items {
		feedItems = append(feedItems, map[string]interface{}{
			"title":       item.Title,
			"link":        item.Link,
			"published":   item.Published,
			"description": item.Description,
		})
	}

	// JSONレスポンスを返却
	return c.JSON(http.StatusOK, map[string]interface{}{
		"title":       feed.Title,
		"description": feed.Description,
		"items":       feedItems,
	})
}

// Repository represents a GitHub repository
type Repository struct {
	Name        string
	Description string
	URL         string
}

// FetchGitHubTrending fetches trending repositories from GitHub
func FetchGitHubTrending() ([]Repository, error) {
	var repositories []Repository

	// Request the HTML page.
	res, err := http.Get("https://github.com/trending")
	if err != nil {
		return repositories, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return repositories, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return repositories, err
	}

	// Find the repository items
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		name := s.Find("h1.h3 a").Text()
		description := s.Find("p.col-9").Text()
		url, _ := s.Find("h1.h3 a").Attr("href")

		repositories = append(repositories, Repository{
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(description),
			URL:         "https://github.com" + strings.TrimSpace(url),
		})
	})

	return repositories, nil
}

// GitHubTrendingHandler handles the /trending endpoint
func GitHubTrendingHandler(c echo.Context) error {
	repos, err := FetchGitHubTrending()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to fetch GitHub trending repositories: %v", err),
		})
	}
	return c.JSON(http.StatusOK, repos)
}
