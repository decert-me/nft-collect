package service

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"nft-collect/internal/app/model"
	"testing"
	"time"
)

func Test_addCollectionByContract(t *testing.T) {
	assetsUrl := fmt.Sprintf("https://polygonapi.nftscan.com/api/v2/account/own/0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266?limit=100&erc_type=erc721&show_attribute=false&sort_field=&sort_direction=")
	client := req.C().SetTimeout(120*time.Second).
		SetCommonRetryCount(1).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("Referer", "https://docs.nftscan.com/").
		SetCommonHeader("X-API-KEY", "xxxxxxx")

	var cursor string
	contractList := make(map[common.Address]struct{})
	contractList[common.HexToAddress("0x60f028c82f9f3bf71e0c13fe9e8e7f916b345c00")] = struct{}{}
	var nft []model.Collection
	for {
		fmt.Println(time.Now())
		reqUrl := assetsUrl + fmt.Sprintf("&cursor=%s", cursor)
		fmt.Println(reqUrl)

		req, errReq := client.R().Get(reqUrl)
		if errReq != nil {
			fmt.Println(errReq)
			break
		}
		res := req.String()

		if gjson.Get(res, "data.total").Uint() == 0 {
			fmt.Println("No data")
			break
		}

		if gjson.Get(res, "code").String() != "200" {
			break
		}
		var nftScan []model.NFTScanOwn
		if errParse := json.Unmarshal([]byte(gjson.Get(res, "data.content").String()), &nftScan); errParse != nil {
			fmt.Println(errParse)
			break
		}
		for _, v := range nftScan {
			if _, ok := contractList[common.HexToAddress(v.ContractAddress)]; !ok {
				continue
			}
			nft = append(nft, model.Collection{Chain: "polygonapi", AccountAddress: "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", NFTScanOwn: v})
		}
		cursor = gjson.Get(res, "data.next").String()
		if cursor == "" {
			break
		}
	}
	fmt.Println(len(nft))
}
