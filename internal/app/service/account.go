package service

import (
	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/utils"
	"nft-collect/pkg/slice"
	"strings"
	"time"
)

// GetContract
// @description:
// @param: address string
// @return: res []response.GetContractRes, err error
func GetContract(address, account string) (res []response.GetContractRes, err error) {
	db := global.DB
	var user model.Account
	var dealList []string
	errFirst := db.Model(&model.Account{}).Where("address = ?", address).First(&user).Error
	if errFirst != nil && errFirst != gorm.ErrRecordNotFound {
		return res, errFirst
	}
	// 获取默认合约
	var contractDefault []string
	err = db.Model(&model.ContractDefault{}).Raw("SELECT contract_id FROM contract_default").Scan(&contractDefault).Error
	if errFirst == gorm.ErrRecordNotFound {
		initAccount(address)
		if err = db.Model(&model.Account{}).Where("address = ?", address).First(&user).Error; err != nil {
			return
		}
		dealList = contractDefault
	} else if address != common.HexToAddress("0").String() {
		dealList = slice.DiffSlice[string](contractDefault, user.ContractIDs)
		go func() {
			time.Sleep(5 * time.Second)
			updateAllCollection(address, dealList, false, false) // update all collection
		}()
	}

	if len(user.ContractIDs) != len(user.Counts) {
		updateContractCount(address)
	}

	contractMap := make(map[string]int64)
	// 查询默认合约
	for _, id := range dealList {
		// TODO 优化
		var count int64
		err = db.Model(&model.Collection{}).
			Raw("SElECT COUNT(1) FROM collection a JOIN contract b ON a.contract_address=b.contract_address AND a.chain=b.chain WHERE b.id = ? AND account_address= ? AND a.status=2", id, address).
			Scan(&count).Error
		contractMap[id] = count
	}
	// show different counts
	if address == account {
		for i, _ := range user.ContractIDs {
			contractMap[user.ContractIDs[i]] = user.Counts[i]
		}
	} else {
		for i, _ := range user.ContractIDs {
			contractMap[user.ContractIDs[i]] = user.CountsShow[i]
		}
	}
	// slice as query
	var contractIDs []string
	for _, c := range dealList {
		contractIDs = append(contractIDs, c)
	}
	for _, c := range user.ContractIDs {
		contractIDs = append(contractIDs, c)
	}
	// get contract detail
	var contract []model.Contract
	err = db.Model(&model.Contract{}).Where("id", contractIDs).Find(&contract).Error
	if err != nil {
		return res, err
	}
	// order by contractIDs
	for _, id := range contractIDs {
		for _, c := range contract {
			if id == c.ID {
				if contractMap[c.ID] == 0 {
					break
				}
				res = append(res, response.GetContractRes{Contract: c, Count: contractMap[c.ID]})
				break
			}
		}
	}
	return res, err
}

// initAccount
// @description: init account
// @param: address string
// @return: err error
func initAccount(address string) (err error) {
	var uuidList []string
	// 创建
	for _, v := range global.CONFIG.NFT.DefContract {
		var id string
		db := global.DB.Model(&model.Contract{}).Select("id")
		temp := strings.Split(v, "::")
		if len(temp) == 0 {
			continue
		}
		db.Where("contract_address", strings.ToLower(temp[0]))
		if len(temp) == 2 {
			db.Where("chain", temp[1])
		}
		errFirst := db.First(&id).Error
		if errFirst != nil {
			if errFirst != gorm.ErrRecordNotFound {
				return errFirst
			}
			continue
		}
		uuidList = append(uuidList, id)
	}
	user := &model.Account{Address: address}
	err = global.DB.Model(&model.Account{}).Create(&user).Error
	if err != nil {
		return err
	}
	updateAllCollection(address, uuidList, true, false)
	return err
}

// updateContractCount
// @description: update the account contract count
// @param: address string
// @return: err error
func updateContractCount(address string) (err error) {
	db := global.DB
	var user model.Account
	if err = db.Model(&model.Account{}).Select("contract_ids").Where("address", address).First(&user).Error; err != nil {
		return err
	}
	var counts []int64
	var countsShow []int64
	var contractIDs []string
	for _, v := range user.ContractIDs {
		var nftContract model.Contract
		if errErrFind := db.Model(&model.Contract{}).Where("id", v).First(&nftContract).Error; errErrFind != nil {
			continue
		}
		var count int64
		var countShow int64

		// TODO: 添加状态 需要确定是否过滤
		err = db.Model(&model.Collection{}).Where("chain", nftContract.Chain).
			Where("contract_address", nftContract.ContractAddress).Where("account_address", address).
			Count(&count).Error
		if err != nil {
			return err
		}
		err = db.Model(&model.Collection{}).Where("chain", nftContract.Chain).
			Where("contract_address", nftContract.ContractAddress).Where("account_address", address).Where("status", 2).
			Count(&countShow).Error
		if err != nil {
			return err
		}
		contractIDs = append(contractIDs, v)
		counts = append(counts, count)
		countsShow = append(countsShow, countShow)
	}
	if err = db.Model(&model.Account{}).Where("address", address).Updates(model.Account{ContractIDs: contractIDs, Counts: counts, CountsShow: countsShow}).Error; err != nil {
		return err
	}

	// TODO: 清除零数量的合约
	return nil
}

// addContractToUser
// @description: add contract to user
// @param: id string, address string
// @return: err error
func addContractToUser(id string, address string) (err error) {
	db := global.DB.Model(&model.Account{})
	var user model.Account
	err = db.Where("address", address).First(&user).Error

	var contractIDs []string
	for _, c := range user.ContractIDs {
		contractIDs = append(contractIDs, c)
	}
	if utils.SliceIsExist[string](contractIDs, id) {
		return
	}

	if err != nil {
		return
	}
	err = db.
		Where("id", user.ID).
		Update("contract_ids", gorm.Expr("\"contract_ids\"||?", "{"+id+"}")).Error

	return err
}
