package session

import "github.com/gin-gonic/gin"

type Handler interface {
	RefreshSession(c *gin.Context)
	DestroySession(c *gin.Context)
}
