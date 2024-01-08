package router

import (
	"github.com/gin-gonic/gin"
	v1 "nft-collect/internal/app/api/v1"
)

func InitZcloakRouter(Router *gin.RouterGroup) {
	router := Router.Group("zcloak")
	{
		router.POST("/saveCardInfo", v1.SaveCardInfo) // SaveCardInfo 保存Zcloak证书
	}
}
