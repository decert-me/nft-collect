package router

import (
	"github.com/gin-gonic/gin"
	"nft-collect/internal/app/api/v1"
	"nft-collect/internal/app/middleware"
)

func InitAccountRouter(Router *gin.RouterGroup) {
	accountRouterWithAddr := Router.Group("account").Use(middleware.Addr())
	accountRouterWithAuth := Router.Group("account").Use(middleware.Auth()) // auth
	{
		accountRouterWithAddr.GET("/own/:address/contract", v1.GetContract)    // Get the list of user NFT contracts
		accountRouterWithAddr.GET("/own/:address", v1.GetCollection)           // Get the NFT data by the user
		accountRouterWithAddr.POST("/own/refreshUserData", v1.RefreshUserData) // refresh user data
	}
	{
		accountRouterWithAuth.GET("/contract/:address", v1.GetCollectionByContract) // Get the NFT data by the Contract
	}
	{
		accountRouterWithAuth.POST("/own/collection", v1.AddCollection)        // add collection
		accountRouterWithAuth.PUT("/own/collection/:id", v1.UpdatedCollection) // update collection status
	}
}
