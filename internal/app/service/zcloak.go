package service

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"strings"
	"time"
)

// SaveCardInfo 保存Zcloak证书
func SaveCardInfo(c *gin.Context, r request.SaveCardInfoRequest) (err error) {
	r.AccountAddress = strings.ToLower(r.AccountAddress)
	// 校验Key
	if c.GetHeader("x-api-key") != global.CONFIG.System.APIKey {
		global.LOG.Error("非法请求", zap.String("x-api-key", c.GetHeader("x-api-key")))
		return errors.New("非法请求")
	}
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
	// 查询NFT是否存在
	var collectionRes model.Collection
	global.DB.
		Select("id,claim_status").
		Where("account_address = ? AND contract_address = ? AND token_id = ?", r.AccountAddress, r.ContractAddress, r.TokenID).
		First(&collectionRes)
	// 已经领取，不需要操作
	if collectionRes.ClaimStatus == 3 {
		return
	} else if collectionRes.ClaimStatus == 2 {
		// 更新数据
		err = global.DB.Model(&model.Collection{}).Where("id = ?", collectionRes.ID).Update("name", r.Name).Error
		if err != nil {
			log.Error("update error", zap.Error(err))
			return
		}
		return
	}
	// 查询NFT是否存在
	var collectionRes2 model.Collection
	externalLink := fmt.Sprintf("https://decert.me/quests/%s", r.TokenID)
	externalLinkUUID := fmt.Sprintf("https://decert.me/quests/%s", r.UUID)
	global.DB.
		Select("id").
		Where("account_address = ? AND contract_name='Decert Badge' AND (external_link = ? OR external_link = ?)", r.AccountAddress, externalLink, externalLinkUUID).
		First(&collectionRes2)
	// 已经领取
	if collectionRes2.ID != "" {
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
			ContractName:    "Decert Badge",
			TokenID:         r.TokenID,
			Owner:           r.AccountAddress,
			Name:            r.Name,
			ErcType:         r.ErcType,
			ImageURI:        r.ImageURI,
			MetadataJSON:    r.MetadataJson,
			MintTimestamp:   time.Now().UnixMilli(),
			ExternalLink:    externalLink,
		},
	}
	err = global.DB.Create(&collection).Error
	if err != nil {
		return
	}
	// 更新合约数量
	updateContractCount(r.AccountAddress)

	return nil
}

// SaveSolanaCardInfo 保存Zcloak证书
func SaveSolanaCardInfo(c *gin.Context, r request.SaveCardInfoRequest) (err error) {
	// 校验Key
	if c.GetHeader("x-api-key") != global.CONFIG.System.APIKey {
		global.LOG.Error("非法请求", zap.String("x-api-key", c.GetHeader("x-api-key")))
		return errors.New("非法请求")
	}
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
	// 查询NFT是否存在
	var collectionRes model.CollectionSolana
	global.DB.
		Select("id,claim_status").
		Where("minter = ? AND collection='Decert Badge' AND token_id = ?", r.AccountAddress, r.TokenID).
		First(&collectionRes)
	// 已经领取，不需要操作
	if collectionRes.ClaimStatus == 3 {
		return
	} else if collectionRes.ClaimStatus == 2 {
		// 更新数据
		err = global.DB.Model(&model.CollectionSolana{}).Where("id = ?", collectionRes.ID).Update("name", r.Name).Error
		if err != nil {
			log.Error("update error", zap.Error(err))
			return
		}
		return
	}
	// 查询NFT是否存在
	var collectionRes2 model.CollectionSolana
	externalLink := fmt.Sprintf("https://decert.me/quests/%s", r.TokenID)
	externalLinkUUID := fmt.Sprintf("https://decert.me/quests/%s", r.UUID)
	global.DB.
		Select("id").
		Where("minter = ? AND collection='Decert Badge' AND (external_link = ? OR external_link = ?)", r.AccountAddress, externalLink, externalLinkUUID).
		First(&collectionRes2)
	// 已经领取
	if collectionRes2.ID != "" {
		return
	}
	if collectionRes2.ID != "" && collectionRes2.ClaimStatus == 1 {
		// 改变状态3
		err = global.DB.Model(&model.CollectionSolana{}).Where("id = ?", collectionRes2.ID).Update("claim_status", 3).Error
		if err != nil {
			return
		}
		// 跳出
		return
	}
	// 写入NFT数据
	collection := model.CollectionSolana{
		Status:      0, // 显示
		ClaimStatus: 2,
		NFTScanSolana: model.NFTScanSolana{
			Collection:    "Decert Badge",
			Minter:        r.AccountAddress,
			Owner:         r.AccountAddress,
			Name:          r.Name,
			ImageURI:      r.ImageURI,
			MetadataJSON:  r.MetadataJson,
			MintTimestamp: time.Now().UnixMilli(),
			ExternalLink:  externalLinkUUID,
			TokenID:       r.TokenID,
		},
	}
	err = global.DB.Create(&collection).Error
	if err != nil {
		return
	}
	// 更新合约数量
	updateContractCount(r.AccountAddress)

	return nil
}
