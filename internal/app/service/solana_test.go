package service

import (
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"nft-collect/internal/app/model"
	"testing"
	"time"
)

func TestSolanaGet(t *testing.T) {
	assetsUrl := fmt.Sprintf("https://solanaapi.nftscan.com/api/sol/assets/collection/Decert Badge?show_attribute=false")
	var cursor string
	var nft []model.CollectionSolana
	var errFlag bool
	client := req.C().SetTimeout(120*time.Second).
		SetCommonRetryCount(1).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("X-API-KEY", "2UEM6Uto")
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
	fmt.Println(errFlag)
	fmt.Println("nft", nft)
}
