package router

import (
	"github.com/gin-gonic/gin"
	v1 "nft-collect/internal/app/api/v1"
)

func InitSystemRouter(Router *gin.RouterGroup) {
	systemRouterWithAuth := Router.Group("system")
	{
		systemRouterWithAuth.GET("/contract/default", v1.GetDefaultContract)
		systemRouterWithAuth.POST("/contract/default", v1.AddDefaultContract)
		systemRouterWithAuth.DELETE("/contract/default", v1.DelDefaultContract)
	}
}
