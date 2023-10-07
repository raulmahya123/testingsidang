package routes

import (
	controller "golangsidang/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) { // membuat routes auth
	incomingRoutes.POST("user/signup", controller.Signup()) // membuat routes signup untuk mengani sigup
	incomingRoutes.POST("user/login", controller.Login())   // membuat routes signin untuk mengani sigin
}
