package routes

import (
	"github.com/aniketDinda/zocket/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(inRoutes *gin.Engine) {

	inRoutes.GET("/health", controllers.Health())
	inRoutes.POST("user/new", controllers.NewUser())
	inRoutes.POST("product/add", controllers.AddProduct())
	inRoutes.GET("product/view", controllers.ViewProducts())
}
