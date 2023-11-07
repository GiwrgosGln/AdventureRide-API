package routes

import (
	"adventureride/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, ctrl *controllers.Controller) {
	r.POST("/register", ctrl.RegisterHandler)
	r.POST("/login", ctrl.LoginHandler)
}
