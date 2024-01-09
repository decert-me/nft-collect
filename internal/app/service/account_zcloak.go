package service

import (
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/response"
)

func GetZcloakContract(didAddress, account string) (res []response.GetContractRes, err error) {
	// 查询用户
	var address string
	err = global.DB.Model(&model.ZcloakDid{}).Select("address").Where("did_address = ?", didAddress).First(&address).Error
	if err != nil {
		return res, nil // 返回空
	}
	// 查询数量
	var count int64
	err = global.DB.Model(&model.Collection{}).Where("account_address = ? AND (claim_status=2 OR claim_status=3)", address).Count(&count).Error
	if err != nil {
		return res, err
	}
	// 添加默认合约
	res = append([]response.GetContractRes{
		{model.Contract{MODEL: global.MODEL{ID: "decert_badge_zcloak"}, ContractName: "Decert Badge", ContractAddress: "", ContractLogo: "ipfs://bafkreiedufqglo2o2shyv2kgvc4hn6j6uajzyrmebi7alowevbmpwhupwy", Chain: ""}, count},
	}, res...)
	return res, err
}
