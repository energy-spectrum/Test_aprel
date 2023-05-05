package route

import (
	"github.com/gin-gonic/gin"

	"app/api/controller"
	"app/config"
	"app/db"
)

func NewAuthRouter(env *config.Env, store db.Store, group *gin.RouterGroup) {
	ac := &controller.AuthController{
		Store: store,
		Env:   env,
	}

	group.POST("/authorize", ac.Authorize)
}
