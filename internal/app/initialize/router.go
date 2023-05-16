package initialize

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"net/http"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/middleware"
	"nft-collect/internal/app/router"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router *gin.Engine
	// 开发环境打开日志 && 打开pprof
	if global.CONFIG.System.Env == "develop" {
		Router = gin.Default()
		pprof.Register(Router) // 性能
	} else {
		Router = gin.New()
		Router.Use(gin.Recovery())
	}
	Router.Use(middleware.Cors())                                                     // 放行跨域请求
	Router.StaticFS(global.CONFIG.NFT.LogoPath, http.Dir(global.CONFIG.NFT.LogoPath)) // 为用户头像和文件提供静态地址
	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
	}
	v1Group := Router.Group("v1")
	{
		router.InitAccountRouter(v1Group)
		router.InitSystemRouter(v1Group)
		router.InitEnsRouter(v1Group)
	}

	global.LOG.Info("router register success")
	return Router
}
