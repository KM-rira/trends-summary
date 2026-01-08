package usecase

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func ScrapeStaticPage(c echo.Context, reqURL string, tags []string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"function": "ScrapeStaticPage",
		"url":      reqURL,
		"tags":     tags,
	}).Info("スクレイピング開始")

	if reqURL == "" {
		logrus.WithFields(logrus.Fields{
			"function":  "ScrapeStaticPage",
			"errorType": "パラメータバリデーションエラー",
		}).Error("URLパラメータが必要です")
		return "", fmt.Errorf("URLパラメータが必要です")
	}

	// URLのバリデーション
	_, err := url.ParseRequestURI(reqURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function":  "ScrapeStaticPage",
			"url":       reqURL,
			"error":     err.Error(),
			"errorType": "URLバリデーションエラー",
		}).Error("URLが無効です")
		return "", fmt.Errorf("URLが無効です: %w", err)
	}

	// タイムアウト付きHTTPクライアントを作成
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// HTTP GETリクエストを送信
	resp, err := client.Get(reqURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function":  "ScrapeStaticPage",
			"url":       reqURL,
			"error":     err.Error(),
			"errorType": "HTTPリクエストエラー",
		}).Error("URLへのリクエストに失敗しました")
		return "", fmt.Errorf("URLへのリクエストに失敗しました: %w", err)
	}
	// プログラム終了時にレスポンスボディを閉じる
	defer resp.Body.Close()

	// ステータスコードが200（OK）か確認
	if resp.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"function":   "ScrapeStaticPage",
			"url":        reqURL,
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
			"errorType":  "HTTPステータスコードエラー",
		}).Error("リクエストが失敗しました")
		return "", fmt.Errorf("リクエストが失敗しました。ステータスコード: %d", resp.StatusCode)
	}

	// goqueryでHTMLを解析
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function":  "ScrapeStaticPage",
			"url":       reqURL,
			"error":     err.Error(),
			"errorType": "HTMLパースエラー",
		}).Error("HTMLの解析に失敗しました")
		return "", fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	getData, err := GetTagDataFromHTML(doc, tags)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function":  "ScrapeStaticPage",
			"url":       reqURL,
			"tags":      tags,
			"error":     err.Error(),
			"errorType": "タグ抽出エラー",
		}).Error("HTMLタグの抽出に失敗しました")
		return "", fmt.Errorf("HTMLの解析に失敗しました: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"function": "ScrapeStaticPage",
		"url":      reqURL,
		"dataLen":  len(getData),
	}).Info("スクレイピング成功")

	return getData, nil
}

func GetTagDataFromHTML(doc *goquery.Document, tags []string) (string, error) {
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

	return getData, nil
}
func GetTagDataFromHTML2(doc *goquery.Document, tags []string) (string, error) {
	var builder strings.Builder
	for i, tag := range tags {
		builder.WriteString(fmt.Sprintf("=== %s #%d ===\n", tag, i+1))

		// 指定されたタグの要素を検索
		doc.Find(tag).Each(func(index int, item *goquery.Selection) {
			// 各行の列データを取得
			item.Find("td").Each(func(j int, cell *goquery.Selection) {
				text := cell.Text()
				builder.WriteString(fmt.Sprintf("Cell %d: %s\n", j+1, strings.TrimSpace(text)))
			})
			builder.WriteString("\n") // 各行ごとの区切り
		})

		builder.WriteString("\n") // タグごとの区切り
	}

	getData := builder.String()
	return getData, nil
}
