package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"trends-summary/internal/usecase"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/google/go-github/github"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Index はRSSフィードを取得して、HTMLまたはJSONで出力するハンドラーです。
func Index(c echo.Context) error {
	logrus.Info("Indexハンドラーが呼び出されました") // デバッグ用ログ

	// RSSフィードURL
	rssURL := "https://feed.infoq.com/"
	logrus.WithField("rssURL", rssURL).Info("RSSフィードURL") // デバッグ用ログ

	// パーサーの作成
	fp := gofeed.NewParser()

	// RSSフィードを取得して解析
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		logrus.WithError(err).Error("RSSフィード取得エラー") // エラーログ
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("RSSフィードの取得に失敗しました: %v", err),
		})
	}

	logrus.WithField("feedTitle", feed.Title).Info("RSSフィード取得成功") // デバッグ用ログ

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
	logrus.Info("GitHubTrendingHandlerが呼び出されました") // デバッグ用ログ

	var trendingRepos []map[string]string

	// GitHub Trendingページをスクレイピング
	res, err := http.Get("https://github.com/trending")
	if err != nil {
		logrus.WithError(err).Error("GitHub Trendingページの取得に失敗しました") // エラーログ
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch GitHub Trending"})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithField("statusCode", res.StatusCode).Error("GitHub Trendingページのステータスコードエラー") // エラーログ
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GitHub Trending page returned non-200 status"})
	}

	// HTMLドキュメントをパース
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.WithError(err).Error("GitHub Trendingページのパースに失敗しました") // エラーログ
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

	logrus.WithField("repositoriesCount", len(trendingRepos)).Info("GitHubトレンドの取得に成功しました") // デバッグ用ログ

	// JSONで返却
	return c.JSON(http.StatusOK, trendingRepos)
}

func TiobeGraph(c echo.Context) error {
	// コンテキストの作成
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// コンテキスト作成
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	// タイムアウト付きのコンテキスト
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second) // タイムアウトを60秒に延長
	defer cancel()

	// ターゲットURL
	url := "https://www.tiobe.com/tiobe-index/"

	// グラフデータを格納する変数
	var graphHTML string

	// タスクを順番に実行
	err := chromedp.Run(ctx, chromedp.Navigate(url))
	if err != nil {
		logrus.Fatalf("Navigation error: %v", err)
	}

	err = chromedp.Run(ctx, chromedp.WaitVisible(`#container`, chromedp.ByID))
	if err != nil {
		logrus.Fatalf("Element visibility error: %v", err)
	}

	err = chromedp.Run(ctx, chromedp.OuterHTML(`#container`, &graphHTML, chromedp.ByID))
	if err != nil {
		logrus.Fatalf("OuterHTML extraction error: %v", err)
	}

	logrus.Info("TIOBEグラフの取得に成功しました") // デバッグ用ログ
	return c.JSON(http.StatusOK, graphHTML)
}

func AIArticleSummary(c echo.Context) error {
	// クエリパラメータからURLを取得
	urlData := c.QueryParam("url")
	if urlData == "" {
		logrus.Fatal("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	htmlData, err := usecase.ScrapeStaticPage(c, urlData, "div.article__data")
	if err != nil {
		logrus.Fatalf("scraping error: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "scraping error"})
	}

	// リクエストURLの構築
	requestText := "下記の記事の内容を簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + htmlData
	logrus.Info("requestText length: ", len(requestText))

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		log.Fatalf("エラー: Gemini APIリクエストに失敗しました: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}

func AIRepositorySummary(c echo.Context) error {
	logrus.Info("AIRepositorySummaryハンドラーが呼び出されました") // デバッグ用ログ
	// クエリパラメータからURLを取得
	urlData := c.QueryParam("url")
	if urlData == "" {
		logrus.Fatal("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	parsedURL, err := url.Parse(urlData)
	if err != nil {
		logrus.Fatal("URL解析エラー:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL解析エラー"})
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		logrus.Fatal("URLが正しい形式ではありません")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが正しい形式ではありません"})
	}
	owner := parts[0]
	repoName := parts[1]

	// コンテキストの作成
	ctx := context.Background()

	// 環境変数から OAuth トークンを取得
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if token == "" {
		log.Fatal("環境変数 GITHUB_OAUTH_TOKEN が設定されていません")
	}

	// OAuth2 クライアントの作成
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// go-github クライアントの作成
	client := github.NewClient(tc)

	// 1. リポジトリ情報の取得
	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		log.Fatalf("リポジトリ情報の取得に失敗しました: %v", err)
	}

	var builder strings.Builder

	fmt.Println("==== リポジトリ情報 ====")
	builder.WriteString(fmt.Sprintf("Name: %s\n", repo.GetName()))
	builder.WriteString(fmt.Sprintf("Description: %s\n", repo.GetDescription()))
	if topics := repo.Topics; len(topics) > 0 {
		fmt.Println("Topics:")
		for _, t := range topics {
			fmt.Printf(" - %s\n", t)
		}
	}

	// 2. README の取得
	readme, _, err := client.Repositories.GetReadme(ctx, owner, repoName, nil)
	if err != nil {
		log.Fatalf("README の取得に失敗しました: %v", err)
	}

	// READMEの内容をデコード（go-github の GetContent メソッドを使用）
	readmeContent, err := readme.GetContent()
	if err != nil {
		log.Fatalf("README のデコードに失敗しました: %v", err)
	}

	builder.WriteString(fmt.Sprintf("==== README 内容 ==== %s", readmeContent))

	repoData := builder.String()

	requestText := "下記はGithubリポジトリのREADMEの内容です。簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + repoData
	logrus.Info("requestText length: ", len(requestText))
	logrus.Info("requestText content: ", requestText)

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		log.Fatalf("エラー: Gemini APIリクエストに失敗しました: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}
func AITrendsSummary(c echo.Context) error {
	// クエリパラメータからURLを取得
	pageData := c.QueryParam("data")
	logrus.Info("pageData length: ", len(pageData))
	if pageData == "" {
		logrus.Fatal("dataパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dataパラメータが必要です。"})
	}

	// pageData を goquery.Document に変換
	reader := strings.NewReader(pageData)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		logrus.Errorf("goquery ドキュメント作成エラー: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました"})
	}
	getData, err := usecase.GetTagDataFromHTML(doc, []string{"table#rss-feed-table"})
	if err != nil {
		log.Fatalf("エラー: HTMLの解析に失敗しました2: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました。2"})
	}
	logrus.Infof("get data: %s", getData)

	requestText := "下記は最新のIT業界のNews一覧です。後述の項目に沿って要約してMarkdown形式で回答してください。・全てのデータから読み取れる傾向と推測される理由 ・InfoQから読み取れる傾向と推測される理由 ・Github daily trendsから読み取れる傾向と推測される理由 ・TIOBE Index Graphから読み取れる傾向と推測される理由\n" + getData
	logrus.Info("requestText length: ", len(requestText))

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		log.Fatalf("エラー: Gemini APIリクエストに失敗しました: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}
