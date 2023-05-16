package service

import (
	"github.com/pkg/errors"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/utils"
	"strings"
)

func GetDefaultContract() (list interface{}, err error) {
	var contractList []model.Contract
	err = global.DB.Raw("SELECT b.* FROM contract_default a LEFT JOIN contract b ON a.contract_id=b.id").Scan(&contractList).Error
	return contractList, err
}

func AddDefaultContract(req request.AddDefaultContractReq) (err error) {
	for _, api := range global.CONFIG.NFT.APIConfig {
		if req.Chain == api.Chain {
			contractMap := map[string]struct{}{strings.ToLower(req.ContractAddress): struct{}{}}
			ItemFiltrateAndDown(contractMap, api)
			break
		}
	}
	var contractID string
	err = global.DB.Raw("SELECT id FROM contract WHERE contract_address=? AND chain = ?", req.ContractAddress, req.Chain).Scan(&contractID).Error
	if err != nil {
		return err
	}
	if contractID == "" {
		return errors.New("contract not exist")
	}
	err = global.DB.Create(&model.ContractDefault{ContractID: contractID}).Error
	return err
}

func DelDefaultContract(req request.DelDefaultContractReq) (err error) {
	var contract model.Contract
	err = global.DB.Model(&model.Contract{}).Where("id", req.ID).First(&contract).Error
	if err != nil {
		return err
	}

	var userList []model.Account
	if err = global.DB.Model(&model.Account{}).Find(&userList).Error; err != nil {
		return
	}
	tx := global.DB.Begin()
	for _, user := range userList {
		if utils.SliceIsExist(user.ContractIDs, req.ID) {
			continue
		}
		err = tx.Model(&model.Collection{}).
			Where("account_address", user.Address).
			Where("contract_address", contract.ContractAddress).
			Where("chain", contract.Chain).
			Delete(&model.Collection{}).Limit(1000).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Model(&model.ContractDefault{}).Where("contract_id", req.ID).Delete(&model.ContractDefault{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}
