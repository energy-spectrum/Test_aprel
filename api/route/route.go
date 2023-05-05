package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"app/config"
	db "app/db"
	"app/api/middleware"
)

func Setup(env *config.Env, store db.Store, router *gin.RouterGroup) {
	router.GET("/ping", func(с *gin.Context) {
		с.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authRouter := router.Group("auth")
	NewAuthRouter(env, store, authRouter)

	authAuditRouter := router.Group("auth-audit")
	authAuditRouter.Use(middleware.TokenAuthMiddleware(store))
	NewAuthAuditRouter(env, store, authAuditRouter)
}
