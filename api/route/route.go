package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/bootstrap"
	db "app/db"
)

func Setup(env *bootstrap.Env, store db.Store, router *gin.RouterGroup) {
	router.GET("/ping", func(с *gin.Context) {
		с.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authRouter := router.Group("auth")
	NewAuthRouter(env, store, authRouter)

	authAuditRouter := router.Group("auth-audit")
	NewAuthAuditRouter(env, store, authAuditRouter)
}
