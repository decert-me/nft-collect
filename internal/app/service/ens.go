package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	ens "github.com/wealdtech/go-ens/v3"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/response"
	"strings"
	"time"
)

// GetEnsRecords 获取ENS记录
func GetEnsRecords(ctx context.Context, q string) (result interface{}, err error) {
	// Get records from Cache
	if strings.Contains(q, ".") {
		var res model.Ens
		errNil := global.DB.Model(&model.Ens{}).Where("domain", q).First(&res).Error
		if errNil != nil {
			return resolveName(ctx, q)
		}
		if res.UpdatedAt.Before(time.Now().AddDate(0, 0, -1)) {
			go resolveName(ctx, q)
		}
		return response.GetEnsResponse{Address: res.Address, Domain: res.Domain, Avatar: res.Avatar}, nil
	} else {
		if !common.IsHexAddress(q) {
			return result, errors.New("RecordNotFound")
		}
		var res model.Ens
		errNil := global.DB.Model(&model.Ens{}).Where("address", q).First(&res).Error
		if errNil != nil {
			return resolveAddress(ctx, common.HexToAddress(q).Hex())
		}
		if res.UpdatedAt.Before(time.Now().AddDate(0, 0, -1)) {
			go resolveAddress(ctx, common.HexToAddress(q).Hex())
		}
		return response.GetEnsResponse{Address: res.Address, Domain: res.Domain, Avatar: res.Avatar}, nil
	}
}

// resolveName 解析 ENS 名称
func resolveName(ctx context.Context, input string) (result response.GetEnsResponse, err error) {
	result.Domain = input
	client, err := ethclient.Dial(global.CONFIG.NFT.EnsRpc)
	if err != nil {
		return result, errors.New("UnexpectedError")
	}
	// Resolve a name to an address.
	address, err := ens.Resolve(client, input)
	if err != nil {
		return result, nil
	}
	result.Address = address.Hex()
	text, errResolver := ens.NewResolver(client, input)
	if errResolver == nil {
		avatar, errAvatar := text.Text("avatar")
		if errAvatar == nil {
			result.Avatar = resolveAvatar(avatar)
		}
	}
	if err != nil {
		return result, errors.New("UnexpectedError")
	}
	ens := model.Ens{Address: address.Hex(), Domain: input, Avatar: result.Avatar}
	err = global.DB.Model(&model.Ens{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}, {Name: "domain"}},
		UpdateAll: true,
	}).Create(&ens).Error
	if err != nil {
		return result, errors.New("UnexpectedError")
	}
	return result, nil
}

// resolveAddress 解析 ENS 地址
func resolveAddress(ctx context.Context, input string) (result response.GetEnsResponse, err error) {
	result.Address = input
	client, err := ethclient.Dial(global.CONFIG.NFT.EnsRpc)
	if err != nil {
		return result, errors.New("UnexpectedError")
	}
	// Resolve address to name
	domain, err := ens.ReverseResolve(client, common.HexToAddress(input))
	if err != nil {
		return result, nil
	}
	result.Domain = domain
	text, errResolver := ens.NewResolver(client, domain)
	if errResolver == nil {
		avatar, errAvatar := text.Text("avatar")
		if errAvatar == nil {
			result.Avatar = resolveAvatar(avatar)
		}
	}

	if err != nil {
		return result, errors.New("UnexpectedError")
	}
	ens := model.Ens{Address: input, Domain: domain, Avatar: result.Avatar}
	if err = global.DB.Model(&model.Ens{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "address"}, {Name: "domain"}},
		UpdateAll: true,
	}).Create(&ens).Error; err != nil {
		return result, errors.New("UnexpectedError")
	}
	return result, nil
}

func resolveAvatar(text string) (res string) {
	if !strings.Contains(text, "eip155") {
		return text
	}
	fmt.Println(text)
	defer func() {
		if err := recover(); err != nil {
			res = text
		}
	}()
	temp := strings.Split(text, "/")
	contractAddr := strings.Split(temp[1], ":")[1]
	assetsUrl := fmt.Sprintf("https://restapi.nftscan.com/api/v2/assets/%s/%s?show_attribute=false", contractAddr, temp[2])
	fmt.Println(assetsUrl)
	client := req.C().SetTimeout(120*time.Second).
		SetCommonRetryCount(1).SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKey).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("Referer", "https://docs.nftscan.com/")

	response, err := client.R().Get(assetsUrl)
	if err != nil {
		global.LOG.Error("Get error ", zap.Error(err))
		return text
	}
	return "ipfs://" + gjson.Get(response.String(), "data.image_uri").String()
}
