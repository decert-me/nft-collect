package initialize

import (
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/service"
	"strings"
)

// 获取默认NFT合约信息
func InitNFTContract() {
	for _, api := range global.CONFIG.NFT.APIConfig {
		for _, v := range global.CONFIG.NFT.DefContract {
			temp := strings.Split(v, "::")
			if len(temp) == 0 || (len(temp) == 2 && temp[1] != api.Chain) {
				continue
			}
			contractMap := map[string]struct{}{strings.ToLower(temp[0]): struct{}{}}
			service.ItemFiltrateAndDown(contractMap, api)
		}
	}
}
