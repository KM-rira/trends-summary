package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

func RequestGemini(c echo.Context, requestText string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"function":       "RequestGemini",
		"requestTextLen": len(requestText),
	}).Info("Gemini APIリクエスト開始")

	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		errmsg := "APIキーが設定されていません。環境変数 GEMINI_API_KEY を設定してください。"
		logrus.WithFields(logrus.Fields{
			"function":  "RequestGemini",
			"errorType": "環境変数エラー",
		}).Error(errmsg)
		return "", fmt.Errorf("%s", errmsg)
	}

	// タイムアウト付きコンテキスト
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// genaiライブラリのデフォルト（v1beta）を使用
	// option.WithAPIKeyで自動的に認証ヘッダーが追加される
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"function":  "RequestGemini",
			"error":     err.Error(),
			"errorType": "Geminiクライアント作成エラー",
		}).Error("Geminiクライアントの作成に失敗しました")
		return "", fmt.Errorf("Geminiクライアントの作成に失敗しました: %w", err)
	}
	defer client.Close()

	modelName := "gemini-2.5-flash"
	logrus.WithFields(logrus.Fields{
		"function":  "RequestGemini",
		"modelName": modelName,
		"endpoint":  "v1beta",
	}).Info("Gemini APIリクエスト送信")

	model := client.GenerativeModel(modelName)
	response, err := model.GenerateContent(ctx, genai.Text(requestText))
	if err != nil {
		// googleapi.Errorの詳細を抽出
		var gerr *googleapi.Error
		if errors.As(err, &gerr) {
			logrus.WithFields(logrus.Fields{
				"function":     "RequestGemini",
				"modelName":    modelName,
				"endpoint":     "v1beta",
				"statusCode":   gerr.Code,
				"errorMessage": gerr.Message,
				"errorBody":    string(gerr.Body),
				"errorType":    "Gemini API通信エラー",
			}).Error("Gemini APIリクエストに失敗しました（詳細）")
		} else {
			logrus.WithFields(logrus.Fields{
				"function":  "RequestGemini",
				"modelName": modelName,
				"endpoint":  "v1beta",
				"error":     err.Error(),
				"errorType": "Gemini API通信エラー",
			}).Error("Gemini APIリクエストに失敗しました")
		}
		return "", fmt.Errorf("Gemini APIリクエストに失敗しました: %w", err)
	}

	// 結果の表示
	summarys := []string{}

	for _, cand := range response.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// 型アサーションと型スイッチを使用して、Textフィールドが存在するか確認し、値を取得する
				switch v := part.(type) {
				case genai.Text:
					summarys = append(summarys, string(v)) // genai.Textはstring型に変換可能
				default:
					logrus.Warn("Part.(type) is not string")
				}
			}
		}
	}
	// summarysの配列データを文字列に全て結合する
	summary := ""
	for _, s := range summarys {
		summary += s
	}

	logrus.WithFields(logrus.Fields{
		"function":   "RequestGemini",
		"summaryLen": len(summary),
		"partsCount": len(summarys),
	}).Info("Gemini API要約生成成功")

	return summary, nil
}
