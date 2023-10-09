package service

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"gorm.io/gorm/clause"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"time"
)

func SolanaGet() (err error) {
	assetsUrl := fmt.Sprintf("https://solanaapi.nftscan.com/api/sol/assets/collection/Decert Badge?show_attribute=false")
	var cursor string
	var nft []model.CollectionSolana
	var errFlag bool
	client := req.C().SetTimeout(120*time.Second).
		SetCommonRetryCount(1).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKey)
	i := 0
	for {
		i++
		reqUrl := assetsUrl + fmt.Sprintf("&cursor=%s", cursor)
		req, errReq := client.R().Get(reqUrl)
		if errReq != nil {
			fmt.Println(errReq)
			errFlag = true
			break
		}
		res := req.String()
		if gjson.Get(res, "data.total").Uint() == 0 {
			fmt.Println("No data")
			break
		}
		if gjson.Get(res, "code").String() != "200" {
			errFlag = true
			break
		}
		var nftScan []model.NFTScanSolana
		if errParse := json.Unmarshal([]byte(gjson.Get(res, "data.content").String()), &nftScan); errParse != nil {
			fmt.Println(errParse)
			errFlag = true
			break
		}
		for _, v := range nftScan {
			nft = append(nft, model.CollectionSolana{NFTScanSolana: v})
		}
		cursor = gjson.Get(res, "data.next").String()
		if cursor == "" {
			break
		}
	}
	if len(nft) == 0 {
		return nil
	}
	// 保存数据
	if err = global.DB.Model(&model.CollectionSolana{}).Omit("status").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "token_address"}},
		UpdateAll: true,
	}).Create(&nft).Error; err != nil {
		return err
	}
	_ = errFlag
	return nil
}
