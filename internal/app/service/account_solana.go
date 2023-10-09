package service

import (
	"gorm.io/gorm"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/response"
)

// GetSolanaContract 获取 Solana 合约
func GetSolanaContract(address, account string) (res []response.GetSolanaContractRes, err error) {
	db := global.DB
	var user model.Account

	errFirst := db.Model(&model.Account{}).Where("address = ?", address).First(&user).Error
	if errFirst != nil && errFirst != gorm.ErrRecordNotFound {
		return res, errFirst
	}
	// 初始化账户
	if errFirst == gorm.ErrRecordNotFound {
		initSolanaAccount(address)
		if err = db.Model(&model.Account{}).Where("address = ?", address).First(&user).Error; err != nil {
			return
		}
	}

	//if len(user.ContractIDs) != len(user.Counts) {
	//	updateSolanaContractCount(address)
	//}
	// TODO 查询默认合约
	dealList := []string{"Decert Badge"}
	contractMap := make(map[string]int64)
	for _, name := range dealList {
		// TODO 优化
		var count int64
		err = db.Model(&model.Collection{}).
			Raw("SElECT COUNT(1) FROM collection_solana a JOIN contract_solana b ON a.collection=b.contract_name WHERE b.contract_name = ? AND minter= ?", name, address).
			Scan(&count).Error
		contractMap[name] = count
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
	var contract []model.ContractSolana
	err = db.Model(&model.ContractSolana{}).Where("contract_name", contractIDs).Find(&contract).Error
	if err != nil {
		return res, err
	}
	res = append(res, response.GetSolanaContractRes{ContractSolana: contract[0], Count: contractMap[contract[0].ContractName]})
	return res, err
}

// initSolanaAccount
func initSolanaAccount(address string) (err error) {
	var uuidList []string
	// 创建
	var id string
	db := global.DB.Model(&model.ContractSolana{}).Select("id")
	errFirst := db.First(&id).Error
	if errFirst != nil {
		if errFirst != gorm.ErrRecordNotFound {
			return errFirst
		}
		return errFirst
	}
	uuidList = append(uuidList, id)

	user := &model.Account{Address: address}
	err = global.DB.Model(&model.Account{}).Create(&user).Error
	if err != nil {
		return err
	}
	//updateSolanaAllCollection(address, uuidList, true, false)
	return err
}

/*
// updateSolanaContractCount
func updateSolanaContractCount(address string) (err error) {
	db := global.DB
	var user model.Account
	if err = db.Model(&model.Account{}).Select("contract_ids").Where("address", address).First(&user).Error; err != nil {
		return err
	}
	var counts []int64
	var countsShow []int64
	var contractIDs []string
	for _, v := range user.ContractIDs {
		var nftContract model.ContractSolana
		if errErrFind := db.Model(&model.ContractSolana{}).Where("id", v).First(&nftContract).Error; errErrFind != nil {
			continue
		}
		var count int64
		var countShow int64

		// TODO: 添加状态 需要确定是否过滤
		err = db.Model(&model.CollectionSolana{}).
			Where("collection", nftContract.ContractName).Where("minter", address).
			Count(&count).Error
		if err != nil {
			return err
		}
		err = db.Model(&model.CollectionSolana{}).
			Where("collection", nftContract.ContractName).Where("minter", address).Where("status", 2).
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
*/
