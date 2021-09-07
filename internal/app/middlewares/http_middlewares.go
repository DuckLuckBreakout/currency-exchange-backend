package middlewares

import (
	"net/http"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/handler_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/jwt_token"

	"github.com/gin-gonic/gin"
)

type MiddlewareManager struct {
	sessUCase session.UseCase
}

func NewMiddlewareManager(sessUCase session.UseCase) *MiddlewareManager {
	return &MiddlewareManager{
		sessUCase: sessUCase,
	}
}

func (mw *MiddlewareManager) TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenAuth, err := jwt_token.ExtractAccessTokenMetadata(c.Request, &configer.AppConfig.Secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, handler_errors.HttpBadUserCredentials)
			return
		}

		userId, err := mw.sessUCase.GetUserIdByAccessToken(tokenAuth.Uuid)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUpgradeRequired, handler_errors.HttpRottenAccessToken)
			return
		}

		c.Set("uuid", tokenAuth.Uuid)
		c.Set("userId", userId)
		c.Next()
	}
}

func (mw *MiddlewareManager) CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, "+
			"X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
