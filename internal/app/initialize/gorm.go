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
		RegisterTables(db)     // 初始化表
		InitSolanaContract(db) // 初始化Solana默认合辑
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
		model.CollectionSolana{},
		model.ContractSolana{},
	)
	if err != nil {
		global.LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	global.LOG.Info("register table success")
}

// InitSolanaContract 初始化Solana默认合辑
func InitSolanaContract(db *gorm.DB) {
	// 判断是否存在
	var count int64
	if err := db.Model(&model.ContractSolana{}).Count(&count).Error; err != nil {
		global.LOG.Error("init Solana Contract failed", zap.Error(err))
		os.Exit(0)
	}
	if count > 0 {
		return
	}
	// 创建默认用户
	contractSolana := model.ContractSolana{
		Chain:        "solana",
		ContractName: "Decert Badge",
		ContractLogo: "ipfs://QmZm3BPXRdwscE5kBNqifHEEDaTJgHy8E77BiTtZ4QkPaR",
		Status:       2,
	}
	if err := db.Create(&contractSolana).Error; err != nil {
		global.LOG.Error("create init user failed", zap.Error(err))
		os.Exit(0)
	}
}
