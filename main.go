package main

import (
	"net/http"
	"os"
	"trends-summary/internal/handlers"

	"github.com/labstack/echo/v4" // バージョンを指定
	"github.com/sirupsen/logrus"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("ログファイルを開けませんでした: %s", err)
	}
	defer logFile.Close()

	// logrusの出力先をログファイルに設定
	logrus.SetOutput(logFile)

	// ログフォーマットを設定（例: JSONフォーマット）
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// ログレベルを設定（例: Infoレベル）
	logrus.SetLevel(logrus.InfoLevel)

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
	e.GET("/github-trending", handlers.GitHubTrendingHandler)
	e.GET("/tiobe-graph", handlers.TiobeGraph)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":9998"))
}
