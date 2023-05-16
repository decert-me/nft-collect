package initialize

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"os"
)

// InitCommonDB 通用数据库
func InitCommonDB() {
	db := GormPgSql("")
	if db != nil {
		global.DB = db
		RegisterTables(db) // 初始化表
	}
}

// RegisterTables 注册数据库表专用
func RegisterTables(db *gorm.DB) {
	err := db.AutoMigrate(
		model.Collection{},
		model.Contract{},
		model.Account{},
		model.ContractDefault{},
		model.Ens{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")
}
