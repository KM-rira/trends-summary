package handlers

import (
	"context"
	"fmt"
	"io"
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
	logrus.WithFields(logrus.Fields{
		"handler": "Index",
		"method":  c.Request().Method,
		"path":    c.Request().URL.Path,
	}).Info("ハンドラー呼び出し")

	// RSSフィードURL
	rssURL := "https://feed.infoq.com"

	// パーサーの作成
	fp := gofeed.NewParser()

	// RSSフィードを取得して解析
	feed, err := fp.ParseURL(rssURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "Index",
			"rssURL":    rssURL,
			"error":     err.Error(),
			"errorType": "RSSフィード取得エラー",
		}).Error("RSSフィードの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("RSSフィードの取得に失敗しました: %v", err),
		})
	}

	logrus.WithFields(logrus.Fields{
		"handler":   "Index",
		"feedTitle": feed.Title,
		"itemCount": len(feed.Items),
	}).Info("RSSフィード取得成功")

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
	targetURL := "https://github.com/trending"
	logrus.WithFields(logrus.Fields{
		"handler":   "GitHubTrendingHandler",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	var trendingRepos []map[string]string

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// GitHub Trendingページをスクレイピング
	res, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GitHubTrendingHandler",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("GitHub Trendingページの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch GitHub Trending"})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"handler":    "GitHubTrendingHandler",
			"targetURL":  targetURL,
			"statusCode": res.StatusCode,
			"status":     res.Status,
			"errorType":  "HTTPステータスコードエラー",
		}).Error("GitHub Trendingページのステータスコードエラー")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GitHub Trending page returned non-200 status"})
	}

	// HTMLドキュメントをパース
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GitHubTrendingHandler",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTMLパースエラー",
		}).Error("GitHub Trendingページのパースに失敗しました")
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

	logrus.WithFields(logrus.Fields{
		"handler":           "GitHubTrendingHandler",
		"repositoriesCount": len(trendingRepos),
	}).Info("GitHubトレンドの取得に成功しました")

	// JSONで返却
	return c.JSON(http.StatusOK, trendingRepos)
}

// GolangRepsitoryTrendingHandler fetches trending repositories from GitHub
func GolangRepsitoryTrendingHandler(c echo.Context) error {
	targetURL := "https://github.com/trending/go"
	logrus.WithFields(logrus.Fields{
		"handler":   "GolangRepsitoryTrendingHandler",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	var trendingRepos []map[string]string

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// GitHub Trendingページをスクレイピング
	res, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GolangRepsitoryTrendingHandler",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("GitHub Trendingページの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch GitHub Trending"})
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"handler":    "GolangRepsitoryTrendingHandler",
			"targetURL":  targetURL,
			"statusCode": res.StatusCode,
			"status":     res.Status,
			"errorType":  "HTTPステータスコードエラー",
		}).Error("GitHub Trendingページのステータスコードエラー")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GitHub Trending page returned non-200 status"})
	}

	// HTMLドキュメントをパース
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GolangRepsitoryTrendingHandler",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTMLパースエラー",
		}).Error("GitHub Trendingページのパースに失敗しました")
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

	logrus.WithFields(logrus.Fields{
		"handler":           "GolangRepsitoryTrendingHandler",
		"repositoriesCount": len(trendingRepos),
	}).Info("Golangトレンドの取得に成功しました")

	// JSONで返却
	return c.JSON(http.StatusOK, trendingRepos)
}

func TiobeGraph(c echo.Context) error {
	targetURL := "https://www.tiobe.com/tiobe-index/"
	logrus.WithFields(logrus.Fields{
		"handler":   "TiobeGraph",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	// chromedpのオプションを設定（headlessモード、Chrome/Chromiumのパス検出）
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	// 実行可能ファイルのパスを検出（google-chrome, chromium, chromium-browserを試行）
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// コンテキストの作成
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// タイムアウト付きのコンテキスト
	ctx, cancel = context.WithTimeout(ctx, 60*time.Second) // タイムアウトを60秒に延長
	defer cancel()

	// グラフデータを格納する変数
	var graphHTML string

	// タスクを順番に実行
	err := chromedp.Run(ctx, chromedp.Navigate(targetURL))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "TiobeGraph",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "ChromeDPナビゲーションエラー",
		}).Error("TIOBEページへのナビゲーションに失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to navigate to TIOBE page"})
	}

	err = chromedp.Run(ctx, chromedp.WaitVisible(`#container`, chromedp.ByID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "TiobeGraph",
			"targetURL": targetURL,
			"selector":  "#container",
			"error":     err.Error(),
			"errorType": "要素表示待機エラー",
		}).Error("TIOBEグラフ要素の表示待機に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load TIOBE graph"})
	}

	err = chromedp.Run(ctx, chromedp.OuterHTML(`#container`, &graphHTML, chromedp.ByID))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "TiobeGraph",
			"targetURL": targetURL,
			"selector":  "#container",
			"error":     err.Error(),
			"errorType": "HTML抽出エラー",
		}).Error("TIOBEグラフのHTML抽出に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to extract TIOBE graph"})
	}

	logrus.WithFields(logrus.Fields{
		"handler":   "TiobeGraph",
		"htmlSize":  len(graphHTML),
		"targetURL": targetURL,
	}).Info("TIOBEグラフの取得に成功しました")
	return c.JSON(http.StatusOK, graphHTML)
}

func AIArticleSummary(c echo.Context) error {
	// クエリパラメータからURLを取得
	urlData := c.QueryParam("url")
	logrus.WithFields(logrus.Fields{
		"handler": "AIArticleSummary",
		"method":  c.Request().Method,
		"path":    c.Request().URL.Path,
		"url":     urlData,
	}).Info("ハンドラー呼び出し")

	if urlData == "" {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIArticleSummary",
			"errorType": "パラメータバリデーションエラー",
		}).Error("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	tags := []string{".article__content"}
	htmlData, err := usecase.ScrapeStaticPage(c, urlData, tags)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIArticleSummary",
			"url":       urlData,
			"tags":      tags,
			"error":     err.Error(),
			"errorType": "スクレイピングエラー",
		}).Error("記事のスクレイピングに失敗しました")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "scraping error"})
	}

	// リクエストURLの構築
	requestText := "下記の記事の内容を簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + htmlData
	logrus.WithFields(logrus.Fields{
		"handler":        "AIArticleSummary",
		"url":            urlData,
		"requestTextLen": len(requestText),
		"scrapedDataLen": len(htmlData),
	}).Info("Gemini APIリクエスト準備完了")

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIArticleSummary",
			"url":       urlData,
			"error":     err.Error(),
			"errorType": "Gemini APIエラー",
		}).Error("Gemini APIリクエストに失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}

func AIRepositorySummary(c echo.Context) error {
	// クエリパラメータからURLを取得
	urlData := c.QueryParam("url")
	logrus.WithFields(logrus.Fields{
		"handler": "AIRepositorySummary",
		"method":  c.Request().Method,
		"path":    c.Request().URL.Path,
		"url":     urlData,
	}).Info("ハンドラー呼び出し")

	if urlData == "" {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"errorType": "パラメータバリデーションエラー",
		}).Error("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	parsedURL, err := url.Parse(urlData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"url":       urlData,
			"error":     err.Error(),
			"errorType": "URL解析エラー",
		}).Error("URLの解析に失敗しました")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL解析エラー"})
	}

	parts := strings.Split(strings.Trim(parsedURL.Path, "/"), "/")
	if len(parts) < 2 {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"url":       urlData,
			"parts":     parts,
			"errorType": "URLフォーマットエラー",
		}).Error("URLが正しい形式ではありません（owner/repo形式が必要）")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが正しい形式ではありません"})
	}
	owner := parts[0]
	repoName := parts[1]

	logrus.WithFields(logrus.Fields{
		"handler":  "AIRepositorySummary",
		"owner":    owner,
		"repoName": repoName,
	}).Info("GitHubリポジトリ情報取得開始")

	// タイムアウト付きコンテキストの作成
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 環境変数から OAuth トークンを取得
	token := os.Getenv("GITHUB_OAUTH_TOKEN")
	if token == "" {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"errorType": "環境変数エラー",
		}).Error("環境変数 GITHUB_OAUTH_TOKEN が設定されていません")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "GitHub OAuth トークンが設定されていません"})
	}

	// OAuth2 クライアントの作成
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	// go-github クライアントの作成
	client := github.NewClient(tc)

	// 1. リポジトリ情報の取得
	repo, _, err := client.Repositories.Get(ctx, owner, repoName)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"owner":     owner,
			"repoName":  repoName,
			"error":     err.Error(),
			"errorType": "GitHub APIエラー",
		}).Error("リポジトリ情報の取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "リポジトリ情報の取得に失敗しました"})
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
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"owner":     owner,
			"repoName":  repoName,
			"error":     err.Error(),
			"errorType": "GitHub APIエラー",
		}).Error("READMEの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "README の取得に失敗しました"})
	}

	// READMEの内容をデコード（go-github の GetContent メソッドを使用）
	readmeContent, err := readme.GetContent()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"owner":     owner,
			"repoName":  repoName,
			"error":     err.Error(),
			"errorType": "READMEデコードエラー",
		}).Error("READMEのデコードに失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "README のデコードに失敗しました"})
	}

	builder.WriteString(fmt.Sprintf("==== README 内容 ==== %s", readmeContent))

	repoData := builder.String()

	requestText := "下記はGithubリポジトリのREADMEの内容です。簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + repoData
	logrus.WithFields(logrus.Fields{
		"handler":        "AIRepositorySummary",
		"owner":          owner,
		"repoName":       repoName,
		"requestTextLen": len(requestText),
		"repoDataLen":    len(repoData),
	}).Info("Gemini APIリクエスト準備完了")

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AIRepositorySummary",
			"owner":     owner,
			"repoName":  repoName,
			"error":     err.Error(),
			"errorType": "Gemini APIエラー",
		}).Error("Gemini APIリクエストに失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	logrus.WithFields(logrus.Fields{
		"handler":    "AIRepositorySummary",
		"owner":      owner,
		"repoName":   repoName,
		"summaryLen": len(summary),
	}).Info("リポジトリ要約生成成功")

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}

// リクエストボディの構造体を定義
type RequestData struct {
	Data string `json:"data"`
}

func AITrendsSummary(c echo.Context) error {
	logrus.WithFields(logrus.Fields{
		"handler": "AITrendsSummary",
		"method":  c.Request().Method,
		"path":    c.Request().URL.Path,
	}).Info("ハンドラー呼び出し")

	// JSONボディを構造体にバインドする
	var reqData RequestData
	if err := c.Bind(&reqData); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AITrendsSummary",
			"error":     err.Error(),
			"errorType": "リクエストバインドエラー",
		}).Error("リクエストデータのバインドに失敗しました")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "リクエストデータが無効です"})
	}
	pageData := reqData.Data

	logrus.WithFields(logrus.Fields{
		"handler":     "AITrendsSummary",
		"pageDataLen": len(pageData),
	}).Info("リクエストデータ受信")

	if pageData == "" {
		logrus.WithFields(logrus.Fields{
			"handler":   "AITrendsSummary",
			"errorType": "パラメータバリデーションエラー",
		}).Error("dataパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "dataパラメータが必要です。"})
	}

	// pageData を goquery.Document に変換
	reader := strings.NewReader(pageData)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AITrendsSummary",
			"error":     err.Error(),
			"errorType": "HTMLパースエラー",
		}).Error("goqueryドキュメント作成に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました"})
	}

	tags := []string{"#rss-feed-table tr", "#github-trending-table tr", "#golangweekly-container"}
	getData, err := usecase.GetTagDataFromHTML(doc, tags)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AITrendsSummary",
			"tags":      tags,
			"error":     err.Error(),
			"errorType": "HTMLタグ抽出エラー",
		}).Error("HTMLタグデータの抽出に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました。2"})
	}

	requestText := "下記は最新のIT業界のNews一覧です。後述の項目に沿って要約してMarkdown形式で回答してください。・全てのデータから読み取れる傾向と推測される理由 ・InfoQから読み取れる傾向と推測される理由 ・Github daily trendsから読み取れる傾向と推測される理由・golangWeeklyから読み取れる傾向と推測される理由\n" + getData
	logrus.WithFields(logrus.Fields{
		"handler":        "AITrendsSummary",
		"requestTextLen": len(requestText),
		"extractedLen":   len(getData),
	}).Info("Gemini APIリクエスト準備完了")

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AITrendsSummary",
			"error":     err.Error(),
			"errorType": "Gemini APIエラー",
		}).Error("Gemini APIリクエストに失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	logrus.WithFields(logrus.Fields{
		"handler":    "AITrendsSummary",
		"summaryLen": len(summary),
	}).Info("トレンド要約生成成功")

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}

func GolangWeeklyContent(c echo.Context) error {
	targetURL := "https://golangweekly.com/rss/"
	logrus.WithFields(logrus.Fields{
		"handler":   "GolangWeeklyContent",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GolangWeeklyContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("GolangWeekly RSSフィードの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch feed"})
	}
	defer resp.Body.Close()

	// Echo の Response を取得
	res := c.Response()

	// 取得した RSS のレスポンスヘッダーをそのまま転送
	for k, values := range resp.Header {
		for _, v := range values {
			res.Header().Add(k, v)
		}
	}
	res.WriteHeader(resp.StatusCode)

	// Body の内容を転送
	if _, err := io.Copy(res, resp.Body); err != nil {
		return err
	}
	return nil
}

func GoogleCloudContent(c echo.Context) error {
	targetURL := "https://cloudblog.withgoogle.com/products/gcp/rss/"
	logrus.WithFields(logrus.Fields{
		"handler":   "GoogleCloudContent",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "GoogleCloudContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("Google Cloud RSSフィードの取得に失敗しました")
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch feed"})
	}
	defer resp.Body.Close()

	// Echo の Response を取得
	res := c.Response()

	// 取得した RSS のレスポンスヘッダーをそのまま転送
	for k, values := range resp.Header {
		for _, v := range values {
			res.Header().Add(k, v)
		}
	}
	res.WriteHeader(resp.StatusCode)

	// Body の内容を転送
	if _, err := io.Copy(res, resp.Body); err != nil {
		return err
	}
	return nil
}
func AWSContent(c echo.Context) error {
	targetURL := "https://aws.amazon.com/blogs/aws/feed/"
	logrus.WithFields(logrus.Fields{
		"handler":   "AWSContent",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// AWSのRSSフィードを取得
	resp, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AWSContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("AWS RSSフィードの取得に失敗しました")
		return c.String(http.StatusInternalServerError, "Error fetching AWS feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"handler":    "AWSContent",
			"targetURL":  targetURL,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
			"errorType":  "HTTPステータスコードエラー",
		}).Error("AWS RSSフィードのステータスコードエラー")
		return c.String(http.StatusInternalServerError, "Error fetching AWS feed: "+resp.Status)
	}

	// レスポンスボディを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AWSContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "レスポンス読み込みエラー",
		}).Error("AWS RSSフィードの読み込みに失敗しました")
		return c.String(http.StatusInternalServerError, "Error reading feed")
	}

	// XMLとして返す（Content-Typeをapplication/xmlに設定）
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	return c.String(http.StatusOK, string(body))
}
func AzureContent(c echo.Context) error {
	targetURL := "https://azure.microsoft.com/ja-jp/blog/feed/"
	logrus.WithFields(logrus.Fields{
		"handler":   "AzureContent",
		"method":    c.Request().Method,
		"path":      c.Request().URL.Path,
		"targetURL": targetURL,
	}).Info("ハンドラー呼び出し")

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// AzureRSSフィードを取得
	resp, err := client.Get(targetURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AzureContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("Azure RSSフィードの取得に失敗しました")
		return c.String(http.StatusInternalServerError, "Error fetching Azure feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"handler":    "AzureContent",
			"targetURL":  targetURL,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
			"errorType":  "HTTPステータスコードエラー",
		}).Error("Azure RSSフィードのステータスコードエラー")
		return c.String(http.StatusInternalServerError, "Error fetching Azure feed: "+resp.Status)
	}

	// レスポンスボディを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler":   "AzureContent",
			"targetURL": targetURL,
			"error":     err.Error(),
			"errorType": "レスポンス読み込みエラー",
		}).Error("Azure RSSフィードの読み込みに失敗しました")
		return c.String(http.StatusInternalServerError, "Error reading feed")
	}

	logrus.WithFields(logrus.Fields{
		"handler":  "AzureContent",
		"bodySize": len(body),
	}).Info("Azure RSSフィード取得成功")

	// XMLとして返す（Content-Typeをapplication/xmlに設定）
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationXML)
	return c.String(http.StatusOK, string(body))
}
