package handler

import (
	"net/http"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/handler_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/usecase_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/jwt_token"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type SessionHandler struct {
	sessUCase session.UseCase
}

func NewHandler(sessUCase session.UseCase) session.Handler {
	return &SessionHandler{
		sessUCase: sessUCase,
	}
}

// RefreshSession godoc
// @Description Обновление токенов access и refresh.
// @Accept json
// @Produce json
// @Param refreshToken header string true "Refresh токен из bearer."
// @Success 200 {object} errors.Error "Обнавление не требуется. Access токен актуален."
// @Success 201 {object} models.Token "Токены успешно обновлены."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректный токен."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /session [PUT]
func (h *SessionHandler) RefreshSession(c *gin.Context) {
	tokenAuth, err := jwt_token.ExtractRefreshTokenMetadata(c.Request, &configer.AppConfig.Secret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, handler_errors.HttpBadUserCredentials)
		return
	}

	newToken, err := h.sessUCase.RefreshSession(tokenAuth.Uuid)
	switch errors.Cause(err) {
	case nil:
		newToken.Sanitize()
		c.JSON(http.StatusCreated, newToken)
		return
	case usecase_errors.UcAccessTokenNotRotten:
		c.JSON(http.StatusOK, handler_errors.HttpSuccess)
		return
	default:
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}
}

// DestroySession godoc
// @Description Выход из профиля. Выданные Access и Refresh токены становятся невалидными.
// @Accept json
// @Produce json
// @Param accessToken header string true "Access токен из bearer."
// @Success 200 {object} errors.Error "Выход из аккаунта успешно выполнен."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректный токен."
// @Failure 426 {object} errors.Error "Access токен протух и требуется обновление."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /session [delete]
func (h *SessionHandler) DestroySession(c *gin.Context) {
	uuid, ok := c.Get("uuid")
	if !ok {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	err := h.sessUCase.DestroySession(uuid.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	c.JSON(http.StatusOK, handler_errors.HttpSuccess)
	return
}
