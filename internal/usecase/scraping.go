package usecase

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func ScrapeStaticPage(c echo.Context, reqURL string, tags ...string) (string, error) {
	if reqURL == "" {
		logrus.Fatal("URLパラメータが必要です")
		return "", c.JSON(http.StatusBadRequest, map[string]string{"error": "URLパラメータが必要です。"})
	}

	// URLのバリデーション
	_, err := url.ParseRequestURI(reqURL)
	if err != nil {
		logrus.Fatalf("URLが無効です。: %v", err)
		return "", c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}

	// HTTP GETリクエストを送信
	resp, err := http.Get(reqURL)
	if err != nil {
		log.Fatalf("エラー: URLへのリクエストに失敗しました: %v", err)
		return "", c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}
	// プログラム終了時にレスポンスボディを閉じる
	defer resp.Body.Close()

	// ステータスコードが200（OK）か確認
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("エラー: リクエストが失敗しました。ステータスコード: %d", resp.StatusCode)
		return "", c.JSON(http.StatusBadRequest, map[string]string{"error": "URLが無効です。"})
	}

	// goqueryでHTMLを解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("エラー: HTMLの解析に失敗しました: %v", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTMLの解析に失敗しました。"})
	}

	var builder strings.Builder
	for i, tag := range tags {
		builder.WriteString(fmt.Sprintf("=== %s #%d ===\n", tag, i+1))

		// 指定されたタグの要素を検索
		doc.Find(tag).Each(func(index int, item *goquery.Selection) {
			// 要素のHTMLを取得
			html, err := item.Html()
			if err != nil {
				log.Printf("エラー: タグ '%s' の要素のHTML取得に失敗しました: %v", tag, err)
				return
			}
			// 取得したHTMLを追加
			builder.WriteString(fmt.Sprintf("=== %s #%d ===\n%s\n\n", tag, index+1, html))
		})

		builder.WriteString("\n") // タグごとの区切り
	}
	getData := builder.String()
	logrus.Infof("get data: %s", getData)

	return getData, nil
}
