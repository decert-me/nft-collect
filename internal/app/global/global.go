package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nft-collect/internal/app/config"
	"nft-collect/pkg/cache"
)

var (
	DB     *gorm.DB             // 数据库链接
	LOG    *zap.Logger          // 日志框架
	CONFIG config.Server        // 配置信息
	Cache  *cache.BigCacheStore //  缓存
)
