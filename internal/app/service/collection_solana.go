package service

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
)

// updateSolanaAllCollection
func updateSolanaAllCollection(address string, uuidList []string, init bool, refresh bool) {
	var err error
	for _, v := range uuidList {
		var contracts string
		if err = global.DB.Model(&model.ContractSolana{}).Select("contract_address").Where("id", v).First(&contracts).Error; err != nil {
			return
		}
		if err = global.DB.Model(&model.CollectionSolana{}).Where("account_address", address).
			Where("contract_address", contracts).
			Where("status=1").
			Updates(map[string]interface{}{"status": 2}).Error; err != nil {
			return
		}
	}
	// 更新用户展示数量
	//_ = updateSolanaContractCount(address)
}

// GetSolanaCollection
func GetSolanaCollection(req request.GetCollectionReq, account string) (total, totalPublic, totalHidden int64, res []response.GetCollection, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.DB.Model(&model.CollectionSolana{}).Select("collection_solana.*,contract.contract_logo").
		Joins("left join contract ON contract.chain=collection_solana.chain AND contract.contract_address=collection.contract_address").
		Where("collection.account_address", req.AccountAddress)
	if req.Search != "" {
		db.Where("token_id ILIKE ? OR name ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.ContractID != "" {
		var contract model.ContractSolana
		if err := global.DB.Model(&model.ContractSolana{}).Where("id", req.ContractID).First(&contract).Error; err != nil {
			global.LOG.Error("error first", zap.Error(err))
			return total, totalPublic, totalHidden, res, err
		}
		db.Where("collection_solana.collection", contract.ContractName)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if err = db.Session(&gorm.Session{}).Where("collection.status", 1).Count(&totalHidden).Error; err != nil {
		return
	}
	if err = db.Session(&gorm.Session{}).Where("collection.status", 2).Count(&totalPublic).Error; err != nil {
		return
	}

	if req.Status != 0 {
		db.Where("collection.status", req.Status)
	} else if req.AccountAddress != account {
		db.Where("collection.status", 2)
	}

	if req.Sort != "asc" && req.Sort != "desc" {
		req.Sort = "desc"
	}
	orderBy := fmt.Sprintf("own_timestamp %s", req.Sort)
	err = db.Limit(limit).Offset(offset).Order(orderBy).Find(&res).Error
	if err != nil {
		return total, totalPublic, totalHidden, res, err
	}

	return
}
