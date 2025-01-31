package usecase

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func RequestGemini(c echo.Context, requestText string) (string, error) {
	// 環境変数からAPIキーを取得
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		errmsg := "APIキーが設定されていません。環境変数 GEMINI_API_KEY を設定してください。"
		logrus.Fatal(errmsg)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": errmsg})
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	response, err := model.GenerateContent(ctx, genai.Text(requestText))
	if err != nil {
		log.Fatal(err)
		return "", c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	logrus.Infof("AI response: %v", response)

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
	logrus.Infof("get summary: %s", summary)

	return summary, nil
}
