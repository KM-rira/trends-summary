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

// GitHubTrendingHandler fetches trending repositories from GitHub
func GitHubTrendingHandler(c echo.Context) error {
	var trendingRepos []map[string]string

	// GitHub Trendingページをスクレイピング
	res, err := http.Get("https://github.com/trending")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch GitHub Trending"})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GitHub Trending page returned non-200 status"})
	}

	// HTMLドキュメントをパース
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse GitHub Trending page"})
	}

	// トレンドリポジトリを抽出
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		// リポジトリ名とURL
		repoAnchor := s.Find("h2.h3 a")
		repoName := strings.TrimSpace(repoAnchor.Text())
		repoName = strings.ReplaceAll(repoName, "\n", " / ") // 改行を " / " に置換
		repoURL, _ := repoAnchor.Attr("href")

		// 説明
		repoDescription := strings.TrimSpace(s.Find("p").Text())

		// 言語
		repoLanguage := strings.TrimSpace(s.Find("[itemprop='programmingLanguage']").Text())

		// スター数
		repoStars := strings.TrimSpace(s.Find("a.Link--muted").First().Text())

		// リポジトリ情報を配列に追加
		trendingRepos = append(trendingRepos, map[string]string{
			"name":        repoName,
			"url":         "https://github.com" + strings.TrimSpace(repoURL),
			"description": repoDescription,
			"language":    repoLanguage,
			"stars":       repoStars,
		})
	})

	// JSONで返却
	return c.JSON(http.StatusOK, trendingRepos)
}
