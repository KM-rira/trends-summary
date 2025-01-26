package main

import (
	"trends-summary/internal/handlers"

	"github.com/labstack/echo/v4" // バージョンを指定
)

func main() {
	e := echo.New()
	// 静的ファイルを提供
	e.Static("/static", "static") // staticディレクトリ内のファイルを提供

	// ルーティングの設定
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html") // HTMLファイルを返す
	})
	// RSSフィード用のエンドポイント
	e.GET("/rss", handlers.Index) // JSONレスポンスを返すエンドポイント
	// GitHubトレンド用のエンドポイント
	e.GET("/trending", handlers.GitHubTrendingHandler)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":9998"))
}
