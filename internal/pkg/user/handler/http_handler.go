package handler

import (
	"net/http"

	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/handler_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/app/errors/usecase_errors"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/models"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/session"
	"github.com/DuckLuckBreakout/currency-exchange-backend/internal/pkg/user"
	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserHandler struct {
	userUCase user.UseCase
	sessUCase session.UseCase
}

func NewHandler(userUCase user.UseCase, sessionUCase session.UseCase) user.Handler {
	return &UserHandler{
		userUCase: userUCase,
		sessUCase: sessionUCase,
	}
}

// SignUp godoc
// @Description Регистрация нового пользователя по email.
// @Accept json
// @Produce json
// @Param user body models.SignupUserRequest true "Данные нового пользователя."
// @Success 201 {object} models.Token "Пользователь успешно создан. Возвращает токены access и refresh."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 409 {object} errors.Error "Пользователь с данным email уже создан."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/auth [post]
func (h *UserHandler) SignUp(c *gin.Context) {
	signupUser := &models.SignupUserRequest{}
	if err := c.Bind(&signupUser); err != nil {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	signupUser.Sanitize()
	err := validator.ValidateStruct(signupUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	userId, err := h.userUCase.SignUp(signupUser)
	if err != nil {
		switch errors.Cause(err) {
		case usecase_errors.UcEmailAlreadyExists:
			c.JSON(http.StatusForbidden, handler_errors.HttpEmailAlreadyExists)
			return
		default:
			c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
			return
		}
	}

	token, err := h.sessUCase.CreateNewSession(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	token.Sanitize()
	c.JSON(http.StatusCreated, token)
	return
}

// GetSelfProfile godoc
// @Description Получение собственного профиля пользователя.
// @Accept json
// @Produce json
// @Param accessToken header string true "Access токен из bearer."
// @Success 200 {object} models.UserData "Профиль пользователя успешно найден."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректный токен."
// @Failure 426 {object} errors.Error "Access токен протух и требуется обновление."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/me [get]
func (h *UserHandler) GetSelfProfile(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	userData, err := h.userUCase.GetSelfProfile(userId.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	userData.Sanitize()
	c.JSON(http.StatusOK, userData)
	return
}

// LogIn godoc
// @Description Авторизация пользователя по email или login.
// @Accept json
// @Produce json
// @Param user body models.LoginUserRequest true "Данные пользователя."
// @Success 200 {object} models.Token "Пользователь успешно авторизован. Выданы Access и Refresh токены."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректные авторизационные данные - неверный пароль."
// @Failure 410 {object} errors.Error "Некорректные авторизационные данные - неверный email или login."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/auth [put]
func (h *UserHandler) LogIn(c *gin.Context) {
	loginUser := &models.LoginUserRequest{}
	if err := c.Bind(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	loginUser.Sanitize()
	err := validator.ValidateStruct(loginUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	userId, err := h.userUCase.LogIn(loginUser)
	if err != nil {
		switch errors.Cause(err) {
		case usecase_errors.UcCanNotFindUser:
			c.JSON(http.StatusGone, handler_errors.HttpBadUserCredentials)
			return
		case usecase_errors.UcPasswordsNotMatch:
			c.JSON(http.StatusUnauthorized, handler_errors.HttpBadUserCredentials)
			return
		default:
			c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
			return
		}
	}

	token, err := h.sessUCase.CreateNewSession(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	token.Sanitize()
	c.JSON(http.StatusOK, token)
	return
}

// DeleteSelfProfile godoc
// @Description Удаление аккаунта пользователя в прииложении. Токены становятся невалидными.
// @Accept json
// @Produce json
// @Param accessToken header string true "Access токен из bearer."
// @Success 200 {object} errors.Error "Аккаунт пользователя успешно удалён."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректный токен."
// @Failure 426 {object} errors.Error "Access токен протух и требуется обновление."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/me [delete]
func (h *UserHandler) DeleteSelfProfile(c *gin.Context) {
	// Delete session (access and refresh tokens)
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

	// Delete profile
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	err = h.userUCase.DeleteSelfProfile(userId.(uint64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	c.JSON(http.StatusOK, handler_errors.HttpSuccess)
	return
}

// CheckUniqEmail godoc
// @Description Выполняется проверка уникальности email.
// @Accept json
// @Produce json
// @Param email query string true "Email который следует проверить."
// @Success 200 {object} errors.Error "Email является уникальным."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 409 {object} errors.Error "Пользователь с таким email уже существует."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/me/email [get]
func (h *UserHandler) CheckUniqEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	emailIsExist, err := h.userUCase.CheckUniqEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	if !emailIsExist {
		c.JSON(http.StatusConflict, handler_errors.HttpEmailAlreadyExists)
		return
	}

	c.JSON(http.StatusOK, handler_errors.HttpSuccess)
	return
}

// UpdateSelfLogin godoc
// @Description Обновление login в профиле пользователя.
// @Accept json
// @Produce json
// @Param login query string true "Login пользователя."
// @Param accessToken header string true "Access токен из bearer."
// @Success 200 {object} errors.Error "Login успешно обновлён."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 409 {object} errors.Error "Пользователь с таким login уже существует."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/me/login [put]
func (h *UserHandler) UpdateSelfLogin(c *gin.Context) {
	username := c.Query("login")
	if username == "" {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	err := h.userUCase.UpdateSelfUsername(userId.(uint64), username)
	if err != nil {
		c.JSON(http.StatusConflict, handler_errors.HttpLoginAlreadyExists)
		return
	}

	c.JSON(http.StatusOK, handler_errors.HttpSuccess)
	return
}

// UpdateSelfAvatar godoc
// @Description Обновление аватарки пользователя. Старая аватарка удаляется.
// @Accept mpfd
// @Produce json
// @Param avatar formData file true "Аватарка пользователя."
// @Param accessToken header string true "Access токен из bearer."
// @Success 200 {object} models.UserAvatar "Аватарка успешно обновлена."
// @Failure 400 {object} errors.Error "Некорректное тело запроса."
// @Failure 401 {object} errors.Error "Некорректный токен."
// @Failure 426 {object} errors.Error "Access токен протух и требуется обновление."
// @Failure 500 {object} errors.Error "Непредвиденная ошибка сервера."
// @Router /user/me/avatar [put]
func (h *UserHandler) UpdateSelfAvatar(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, handler_errors.HttpIncorrectRequestBody)
		return
	}

	avatar, err := h.userUCase.UpdateSelfAvatar(userId.(uint64), fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, handler_errors.HttpInternalServerError)
		return
	}

	avatar.Sanitize()
	c.JSON(http.StatusOK, avatar)
	return
}
