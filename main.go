package main

import (
	"net/http"
	"os"
	"time"
	"trends-summary/internal/handlers"
	"trends-summary/internal/middleware"

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

	// 認証不要のエンドポイント
	e.POST("/api/login", handlers.Login)
	e.POST("/api/logout", handlers.Logout)
	
	// 認証が必要なAPIグループ
	api := e.Group("")
	api.Use(middleware.AuthMiddleware)
	
	// 認証状態確認
	api.GET("/api/check-auth", handlers.CheckAuth)

	// RSSフィード用のエンドポイント（英語版）
	api.GET("/rss", handlers.Index)      // JSONレスポンスを返すエンドポイント
	api.GET("/rss-ja", handlers.IndexJA) // 日本語版

	// GitHubトレンド用のエンドポイント
	api.GET("/github-trending", handlers.GitHubTrendingHandler)
	api.GET("/golang-repository-trending", handlers.GolangRepsitoryTrendingHandler)
	api.GET("/tiobe-graph", handlers.TiobeGraph)
	api.GET("/ai-article-summary", handlers.AIArticleSummary)
	api.GET("/ai-repository-summary", handlers.AIRepositorySummary)
	api.GET("/golang-weekly-content", handlers.GolangWeeklyContent)

	// クラウドRSSフィード（英語版）
	api.GET("/google-cloud-content", handlers.GoogleCloudContent)
	api.GET("/aws-content", handlers.AWSContent)
	api.GET("/azure-content", handlers.AzureContent)

	// クラウドRSSフィード（日本語版）
	api.GET("/google-cloud-content-ja", handlers.GoogleCloudContentJA)
	api.GET("/aws-content-ja", handlers.AWSContentJA)
	api.GET("/azure-content-ja", handlers.AzureContentJA)

	api.POST("/ai-trends-summary", handlers.AITrendsSummary)

	// 静的ファイルを提供（ワイルドカードの前に配置することが重要）
	e.Static("/assets", "static/assets")
	e.Static("/static", "static")
	e.File("/vite.svg", "static/vite.svg")

	// ルートアクセスでindex.htmlを返す
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}
