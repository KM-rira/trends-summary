package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// loggingTransport は HTTPリクエストをログ出力するためのトランスポート
type loggingTransport struct {
	wrapped http.RoundTripper
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// リクエストURLをログ出力
	logrus.WithFields(logrus.Fields{
		"function": "RequestGemini",
		"method":   req.Method,
		"url":      req.URL.String(),
		"host":     req.URL.Host,
		"path":     req.URL.Path,
	}).Info("実際のHTTPリクエスト")

	// 実際のリクエストを実行
	resp, err := t.wrapped.RoundTrip(req)

	if err == nil {
		logrus.WithFields(logrus.Fields{
			"function":   "RequestGemini",
			"statusCode": resp.StatusCode,
			"status":     resp.Status,
		}).Info("HTTPレスポンス")
	}

	return resp, err
}

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

	// HTTPトランスポートをラップしてリクエストURLをログ出力
	httpClient := &http.Client{
		Transport: &loggingTransport{
			wrapped: http.DefaultTransport,
		},
	}

	// genaiライブラリのデフォルトはv1betaなので、エンドポイント指定は不要
	// option.WithHTTPClientでログ出力するクライアントを使用
	client, err := genai.NewClient(ctx,
		option.WithAPIKey(apiKey),
		option.WithHTTPClient(httpClient),
	)
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
