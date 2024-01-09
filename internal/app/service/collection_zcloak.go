package service

import (
	"fmt"
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
)

func GetZcloakCollection(req request.GetCollectionReq, account string) (total, totalPublic, totalHidden int64, res []response.GetCollection, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.DB.Model(&model.Collection{}).Select("collection.*,contract.contract_logo").
		Joins("left join contract ON contract.chain=collection.chain AND contract.contract_address=collection.contract_address").
		Joins("left join zcloak_did ON zcloak_did.address=collection.account_address").
		Where("zcloak_did.did_address", req.AccountAddress).Where("collection.claim_status = 2 OR collection.claim_status = 3")
	if req.Search != "" {
		db.Where("token_id ILIKE ? OR name ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	//if req.ContractID != "" {
	//	var contract model.Contract
	//	if err := global.DB.Model(&model.Contract{}).Where("id", req.ContractID).First(&contract).Error; err != nil {
	//		global.LOG.Error("error first", zap.Error(err))
	//		return total, totalPublic, totalHidden, res, err
	//	}
	//	db.Where("collection.chain", contract.Chain).Where("collection.contract_address", contract.ContractAddress)
	//}
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
		req.Sort = "asc"
	}
	orderBy := fmt.Sprintf("mint_timestamp %s", req.Sort)
	err = db.Limit(limit).Offset(offset).Order(orderBy).Find(&res).Error
	if err != nil {
		return total, totalPublic, totalHidden, res, err
	}

	return
}
