package router

import (
	"github.com/gin-gonic/gin"
	"nft-collect/internal/app/api/v1"
	"nft-collect/internal/app/middleware"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	accountRouterWithCache := Router.Group("account").Use(middleware.Addr())
	accountRouterWithAuth := Router.Group("account").Use(middleware.Auth()) // auth
	{
		accountRouterWithCache.GET("/own/:address/contract", v1.GetContract) // Get the list of user NFT contracts
		accountRouterWithCache.GET("/own/:address", v1.GetCollection)        // Get the NFT data by the user
	}
	{
		accountRouterWithAuth.GET("/contract/:address", v1.GetCollectionByContract) // Get the NFT data by the Contract
	}
	{
		accountRouterWithAuth.POST("/own/collection", v1.AddCollection)        // add collection
		accountRouterWithAuth.PUT("/own/collection/:id", v1.UpdatedCollection) // update collection status
		accountRouterWithAuth.POST("/own/refreshUserData", v1.RefreshUserData) // add collection

	}
}
