package initialize

import (
	"nft-collect/internal/app/global"
	"nft-collect/pkg/cache"
	"time"
)

func InitCache() {
	if global.CONFIG.NFT.CacheTime < 3 {
		panic("Less than 3 minutes")
	}
	global.Cache = cache.NewBigCacheStore(time.Duration(global.CONFIG.NFT.CacheTime)*time.Minute, global.LOG)
}
