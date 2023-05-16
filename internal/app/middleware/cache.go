package middleware

import (
	cache "github.com/chenyahui/gin-cache"
	"github.com/gin-gonic/gin"
	"nft-collect/internal/app/global"
)

func Cache() gin.HandlerFunc {
	return cache.CacheByRequestURI(global.Cache, 0)
}
