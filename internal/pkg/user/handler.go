package user

import "github.com/gin-gonic/gin"

type Handler interface {
	SignUp(c *gin.Context)
	LogIn(c *gin.Context)
	GetSelfProfile(c *gin.Context)
	DeleteSelfProfile(c *gin.Context)
	CheckUniqEmail(c *gin.Context)
	UpdateSelfLogin(c *gin.Context)
	UpdateSelfAvatar(c *gin.Context)
}
