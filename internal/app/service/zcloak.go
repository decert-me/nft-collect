package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"strings"
)

// SaveCardInfo 保存Zcloak证书
func SaveCardInfo(c *gin.Context, r request.SaveCardInfoRequest) (err error) {
	// 校验Key
	if c.GetHeader("x-api-key") != global.CONFIG.System.APIKey {
		global.LOG.Error("非法请求", zap.String("x-api-key", c.GetHeader("x-api-key")))
		return errors.New("非法请求")
	}
	// 查询NFT是否存在
	var collectionRes model.Collection
	global.DB.
		Select("id,claim_status").
		Where("account_address = ? AND contract_address = ? AND token_id = ?", r.AccountAddress, r.ContractAddress, r.TokenID).
		First(&collectionRes)
	// 已经领取，不需要操作
	if collectionRes.ClaimStatus == 3 || collectionRes.ClaimStatus == 2 {
		return
	}
	if collectionRes.ID != "" && collectionRes.ClaimStatus == 1 {
		// 改变状态3
		err = global.DB.Model(&model.Collection{}).Where("id = ?", collectionRes.ID).Update("claim_status", 3).Error
		if err != nil {
			return
		}
		// 跳出
		return
	}
	// 写入NFT数据
	collection := model.Collection{
		Chain:          r.Chain,
		AccountAddress: r.AccountAddress,
		Status:         2, // 显示
		ClaimStatus:    2,
		NFTScanOwn: model.NFTScanOwn{
			ContractAddress: r.ContractAddress,
			TokenID:         r.TokenID,
			Owner:           r.AccountAddress,
			Name:            r.Name,
			ErcType:         r.ErcType,
			ImageURI:        r.ImageURI,
			MetadataJSON:    r.MetadataJson,
		},
	}
	err = global.DB.Create(&collection).Error
	if err != nil {
		return
	}
	// 更新合约数量
	updateContractCount(r.AccountAddress)
	// 保存did和地址映射
	var count int64
	err = global.DB.Model(&model.ZcloakDid{}).Where("address = ?", r.AccountAddress).Count(&count).Error
	if err != nil {
		return
	}
	if count > 0 {
		return
	}
	// 插入
	err = global.DB.Create(&model.ZcloakDid{
		Address:    strings.ToLower(r.AccountAddress),
		DidAddress: r.DidAddress,
	}).Error
	if err != nil {
		return
	}
	return nil
}
