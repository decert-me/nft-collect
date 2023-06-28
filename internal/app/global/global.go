package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nft-collect/internal/app/config"
)

var (
	DB        *gorm.DB        // 数据库链接
	LOG       *zap.Logger     // 日志框架
	CONFIG    config.Server   // 配置信息
	ChainName map[uint]string // 链名称
)
