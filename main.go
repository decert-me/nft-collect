package main

import (
	"go.uber.org/zap"
	"nft-collect/internal/app/core"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/initialize"
)

func main() {
	// 初始化Viper
	core.Viper()
	// 初始化zap日志库
	global.LOG = core.Zap()
	// 注册全局logger
	zap.ReplaceGlobals(global.LOG)
	// 初始化数据库
	initialize.InitCommonDB()
	// 初始化默认合约
	initialize.InitNFTContract()
	// 初始化链名称
	initialize.InitChainName()
	core.RunWindowsServer()
}
