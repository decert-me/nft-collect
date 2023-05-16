package service

import (
	"fmt"
	"nft-collect/internal/app/global"
	"strings"
)

func DeleteCache(address string) {
	iterator := global.Cache.Cache.Iterator()
	url := fmt.Sprintf("/account/own/%s", address)
	for iterator.SetNext() {
		entryInfo, err := iterator.Value()
		if err == nil {
			// 处理key和value
			if strings.Contains(strings.ToLower(entryInfo.Key()), url) {
				fmt.Println("delete:", entryInfo.Key())
				global.Cache.Cache.Delete(entryInfo.Key())
			}
		}
	}
}
