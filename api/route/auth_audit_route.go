package route

import (
	"github.com/gin-gonic/gin"

	"app/api/controller"
	"app/config"
	"app/db"
)

func NewAuthAuditRouter(env *config.Env, store db.Store, group *gin.RouterGroup) {
	ac := &controller.AuthAuditController{
		Store: store,
		Env:   env,
	}

	group.GET("/history", ac.GetAuthAudit)
	group.DELETE("/clear", ac.ClearAuthAudit)
}
