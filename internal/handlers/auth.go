package handlers

import (
	"net/http"
	"os"
	"time"

	"trends-summary/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// LoginRequest ログインリクエスト
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login ログインハンドラー
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエストが不正です",
		})
	}

	// 環境変数から認証情報を取得（必須）
	validUsername := os.Getenv("AUTH_USERNAME")
	validPassword := os.Getenv("AUTH_PASSWORD")

	if validUsername == "" || validPassword == "" {
		logrus.Error("AUTH_USERNAME または AUTH_PASSWORD 環境変数が設定されていません")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "サーバー設定エラー",
		})
	}

	// 認証チェック
	if req.Username != validUsername || req.Password != validPassword {
		logrus.WithFields(logrus.Fields{
			"username": req.Username,
		}).Warn("ログイン失敗")
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "ユーザー名またはパスワードが正しくありません",
		})
	}

	// JWTトークン生成
	token, err := middleware.GenerateToken(req.Username)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("JWT生成エラー")
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "認証トークンの生成に失敗しました",
		})
	}

	// Cookieにトークンをセット（24時間有効）
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true // HTTPSのみ
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	logrus.WithFields(logrus.Fields{
		"username": req.Username,
	}).Info("ログイン成功")

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ログインに成功しました",
	})
}

// Logout ログアウトハンドラー
func Logout(c echo.Context) error {
	// Cookieを削除
	cookie := new(http.Cookie)
	cookie.Name = "auth_token"
	cookie.Value = ""
	cookie.Expires = time.Now().Add(-1 * time.Hour) // 過去の時刻を設定
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteStrictMode
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ログアウトしました",
	})
}

// CheckAuth 認証状態確認ハンドラー
func CheckAuth(c echo.Context) error {
	username := c.Get("username")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"authenticated": true,
		"username":      username,
	})
}
