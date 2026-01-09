package main

import (
	"net/http"
	"os"
	"time"
	"trends-summary/internal/handlers"

	"github.com/labstack/echo/v4" // バージョンを指定
	"github.com/sirupsen/logrus"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// logrusの出力先を標準出力に設定（journaldが収集）
	logrus.SetOutput(os.Stdout)

	// ログフォーマットを設定（JSONフォーマット - journaldと相性が良い）
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// ログレベルを設定（例: Infoレベル）
	logrus.SetLevel(logrus.InfoLevel)

	e := echo.New()

	// サーバータイムアウトの設定
	e.Server.ReadTimeout = 30 * time.Second
	e.Server.WriteTimeout = 60 * time.Second
	e.Server.IdleTimeout = 120 * time.Second

	// RSSフィード用のエンドポイント（英語版）
	e.GET("/rss", handlers.Index) // JSONレスポンスを返すエンドポイント
	e.GET("/rss-ja", handlers.IndexJA) // 日本語版
	
	// GitHubトレンド用のエンドポイント
	e.GET("/github-trending", handlers.GitHubTrendingHandler)
	e.GET("/golang-repository-trending", handlers.GolangRepsitoryTrendingHandler)
	e.GET("/tiobe-graph", handlers.TiobeGraph)
	e.GET("/ai-article-summary", handlers.AIArticleSummary)
	e.GET("/ai-repository-summary", handlers.AIRepositorySummary)
	e.GET("/golang-weekly-content", handlers.GolangWeeklyContent)
	
	// クラウドRSSフィード（英語版）
	e.GET("/google-cloud-content", handlers.GoogleCloudContent)
	e.GET("/aws-content", handlers.AWSContent)
	e.GET("/azure-content", handlers.AzureContent)
	
	// クラウドRSSフィード（日本語版）
	e.GET("/google-cloud-content-ja", handlers.GoogleCloudContentJA)
	e.GET("/aws-content-ja", handlers.AWSContentJA)
	e.GET("/azure-content-ja", handlers.AzureContentJA)
	
	e.POST("/ai-trends-summary", handlers.AITrendsSummary)

	// 静的ファイルを提供（ワイルドカードの前に配置することが重要）
	e.Static("/assets", "static/assets")
	e.Static("/static", "static")
	e.File("/vite.svg", "static/vite.svg")

	// ルートアクセスでindex.htmlを返す
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(":80"))
}
