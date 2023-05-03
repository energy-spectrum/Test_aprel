package route

import (
	"github.com/gin-gonic/gin"

	"app/api/controller"
	"app/bootstrap"
	"app/db"
)

func NewAuthAuditRouter(env *bootstrap.Env, store db.Store, group *gin.RouterGroup) {
	ac := &controller.AuthAuditController{
		Store: store,
		Env:   env,
	}

	group.GET("/history", ac.GetAuditEvents)
	group.DELETE("/clear", ac.ClearAudit)
}
