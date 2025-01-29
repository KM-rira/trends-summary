package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"trends-summary/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func RequestGemini(c echo.Context, requestText string) (string, error) {
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logrus.Fatal("APIキーが設定されていません。環境変数 GEMINI_API_KEY を設定してください。")
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "APIキーが設定されていません。環境変数 GEMINI_API_KEY を設定してください。"})
	}
	aiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	// リクエストボディの作成
	requestBody := utils.GenerateContentRequest{
		Contents: []utils.Content{
			{
				Parts: []utils.ContentPart{
					{
						Text: requestText,
					},
				},
			},
		},
	}
	// JSONエンコード
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logrus.Error("リクエストボディのJSONエンコードに失敗:", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "リクエストボディのJSONエンコードに失敗しました。"})
	}

	// HTTP POSTリクエストの作成
	req, err := http.NewRequest("POST", aiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.Error("HTTPリクエストの作成に失敗:", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTTPリクエストの作成に失敗しました。"})
	}

	// ヘッダーの設定
	req.Header.Set("Content-Type", "application/json")

	// HTTPクライアントの作成（タイムアウト設定）
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// リクエストの送信
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("HTTPリクエストの送信に失敗:", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "HTTPリクエストの送信に失敗しました。"})
	}
	defer resp.Body.Close()

	// ステータスコードの確認
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logrus.Errorf("リクエストが失敗しました。ステータスコード: %d\nレスポンスボディ: %s\n", resp.StatusCode, string(bodyBytes))
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "リクエストが失敗しました。"})
	}

	// レスポンスボディの読み取り
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("レスポンスボディの読み取りに失敗:", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "レスポンスボディの読み取りに失敗しました。"})
	}

	// レスポンスJSONのパース
	var response utils.GenerateContentResponse
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		logrus.Error("レスポンスJSONのパースに失敗:", err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "レスポンスJSONのパースに失敗しました。"})
	}
	logrus.Infof("AI response: %v", response)

	// 結果の表示
	summary := ""
	if len(response.Candidates) > 0 {
		candidate := response.Candidates[0]
		if len(candidate.Content.Parts) > 0 {
			logrus.Infof("生成されたコンテンツ: %s", candidate.Content.Parts[0].Text)
			summary = candidate.Content.Parts[0].Text
		} else {
			logrus.Fatal("コンテンツのパーツがありません。")
			return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "コンテンツのパーツがありません。"})
		}
	} else {
		logrus.Fatal("生成されたコンテンツがありません。")
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": "生成されたコンテンツがありません。"})
	}

	return summary, nil
}
