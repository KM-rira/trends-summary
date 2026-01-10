package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Fatal("JWT_SECRET 環境変数が設定されていません。サーバーを起動できません。")
	}
	jwtSecret = []byte(secret)
}

// JWTClaims カスタムクレーム
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken JWT生成
func GenerateToken(username string) (string, error) {
	claims := &JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware 認証ミドルウェア
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Cookieからトークンを取得
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("認証Cookie取得失敗")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "認証が必要です",
			})
		}

		// トークンを検証
		token, err := jwt.ParseWithClaims(cookie.Value, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Warn("JWT検証失敗")
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "無効な認証トークンです",
			})
		}

		// クレームを取得してContextに保存
		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("username", claims.Username)
		}

		return next(c)
	}
}
