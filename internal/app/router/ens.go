package router

import (
	"github.com/gin-gonic/gin"
	v1 "nft-collect/internal/app/api/v1"
)

func InitEnsRouter(Router *gin.RouterGroup) {
	ensRouter := Router.Group("ens")
	{
		ensRouter.GET("/:q", v1.GetEnsRecords)
	}
}
