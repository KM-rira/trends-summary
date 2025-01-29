package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"trends-summary/internal/usecase"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
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
	articleURL := c.QueryParam("url")
	if articleURL == "" {
		logrus.Fatal("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	// URLのバリデーション
	_, err := url.ParseRequestURI(articleURL)
	if err != nil {
		logrus.Fatalf("URLが無効です。: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}

	// HTTP GETリクエストを送信
	resp, err := http.Get(articleURL)
	if err != nil {
		log.Fatalf("エラー: URLへのリクエストに失敗しました: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}
	// プログラム終了時にレスポンスボディを閉じる
	defer resp.Body.Close()

	// ステータスコードが200（OK）か確認
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("エラー: リクエストが失敗しました。ステータスコード: %d", resp.StatusCode)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}

	// goqueryでHTMLを解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("エラー: HTMLの解析に失敗しました: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました。"})
	}

	// <div class="article__data">要素を取得
	aritcleData := ""
	doc.Find("div.article__data").Each(func(index int, item *goquery.Selection) {
		// 選択した要素のHTMLを取得
		html, err := item.Html()
		if err != nil {
			log.Printf("エラー: 要素のHTML取得に失敗しました: %v", err)
		}
		// 取得したHTMLを表示
		aritcleData = fmt.Sprintf("=== div.article__data #%d ===\n%s\n\n", index+1, html)
	})

	// リクエストURLの構築
	requestText := "下記の記事の内容を簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + aritcleData
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
	// クエリパラメータからURLを取得
	articleURL := c.QueryParam("url")
	if articleURL == "" {
		logrus.Fatal("URLパラメータが必要です")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	// URLのバリデーション
	_, err := url.ParseRequestURI(articleURL)
	if err != nil {
		logrus.Fatalf("URLが無効です。: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}

	// HTTP GETリクエストを送信
	resp, err := http.Get(articleURL)
	if err != nil {
		log.Fatalf("エラー: URLへのリクエストに失敗しました: %v", err)
	}
	// プログラム終了時にレスポンスボディを閉じる
	defer resp.Body.Close()

	// ステータスコードが200（OK）か確認
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("エラー: リクエストが失敗しました。ステータスコード: %d", resp.StatusCode)
	}

	// goqueryでHTMLを解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("エラー: HTMLの解析に失敗しました: %v", err)
	}

	// GitHubのREADMEを取得
	articleData := ""
	readmeElement := doc.Find("article.markdown-body.entry-content.container-lg")
	if readmeElement.Length() > 0 {
		html, err := readmeElement.Html()
		if err != nil {
			log.Printf("エラー: 要素のHTML取得に失敗しました: %v", err)
		} else {
			articleData = fmt.Sprintf("=== GitHub README ===\n%s\n\n", html)
		}
	} else {
		log.Println("エラー: README要素が見つかりませんでした")
	}

	requestText := "下記はGithubリポジトリのREADMEの内容です。簡潔に日本語で要約してください。その際に、結果から記載してください。\n" + articleData
	logrus.Info("requestText length: ", len(requestText))

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

	requestText := "下記は最新のIT業界のNews一覧です。後述の項目に沿って要約して回答してください。・全てのデータから読み取れる傾向と推測される理由 ・InfoQから読み取れる傾向と推測される理由 ・Github daily trendsから読み取れる傾向と推測される理由 ・TIOBE Index Graphから読み取れる傾向と推測される理由\n" + pageData
	logrus.Info("requestText length: ", len(requestText))

	summary, err := usecase.RequestGemini(c, requestText)
	if err != nil {
		log.Fatalf("エラー: Gemini APIリクエストに失敗しました: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gemini APIリクエストに失敗しました。"})
	}

	// JSONオブジェクトとしてサマリーを返す
	return c.JSON(http.StatusOK, map[string]string{"summary": summary})
}
