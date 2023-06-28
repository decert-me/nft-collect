package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imroc/req/v3"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nft-collect/internal/app/config"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
	"nft-collect/pkg/slice"
	"sync"
	"time"
)

// UpdatedCollection
// @description: updated Collection status
// @param: req request.GetCollectionReq
// @return: total int64, res []model.NFT, err error
func UpdatedCollection(req request.UpdatedCollectionReq, address string) (err error) {
	raw := global.DB.Model(&model.Collection{}).Where("id = ?", req.ID).Update("status", req.Status)
	if raw.RowsAffected == 0 {
		return errors.New("error")
	}
	// 插入配置
	var contractID string
	err = global.DB.Raw("SELECT b.id FROM collection a LEFT JOIN contract b ON a.chain=b.chain AND a.contract_address=b.contract_address WHERE a.id= ? ", req.ID).
		Scan(&contractID).Error
	if err != nil {
		return err
	}
	addContractToUser(contractID, address)
	go updateContractCount(address)
	return raw.Error
}

// GetCollection
// @description: get collection
// @param: req request.GetCollectionReq
// @return: total int64, res []model.NFT, err error
func GetCollection(req request.GetCollectionReq, account string) (total, totalPublic, totalHidden int64, res []response.GetCollection, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.DB.Model(&model.Collection{}).Select("collection.*,contract.contract_logo").
		Joins("left join contract ON contract.chain=collection.chain AND contract.contract_address=collection.contract_address").
		Where("collection.account_address", req.AccountAddress)
	if req.Search != "" {
		db.Where("token_id ILIKE ? OR name ILIKE ?", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	if req.ContractID != "" {
		var contract model.Contract
		if err := global.DB.Model(&model.Contract{}).Where("id", req.ContractID).First(&contract).Error; err != nil {
			global.LOG.Error("error first", zap.Error(err))
			return total, totalPublic, totalHidden, res, err
		}
		db.Where("collection.chain", contract.Chain).Where("collection.contract_address", contract.ContractAddress)
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

// GetCollectionByContract
// @description: get collection by contract
// @param: req request.GetCollectionReq
// @return: total int64, res []model.NFT, err error
func GetCollectionByContract(req request.GetCollectionReq) (total int64, res []model.Collection, err error) {
	var api config.APIConfig
	for _, v := range global.CONFIG.NFT.APIConfig {
		if v.ChainID == req.ChainID {
			api = v
			break
		}
	}
	if api.Chain == "" {
		return total, res, errors.New("not support chain")
	}
	if req.Page == 1 {
		var firstCollection model.Collection
		errFind := global.DB.Model(&model.Collection{}).
			Select("updated_at").
			Where("account_address", req.AccountAddress).
			Where("contract_address", req.ContractAddress).
			Order("updated_at desc").
			First(&firstCollection).Error
		if errFind == gorm.ErrRecordNotFound {
			wg := new(sync.WaitGroup)
			wg.Add(2)
			go addCollectionByContract(wg, req.AccountAddress, "erc721", api, req.ContractAddress)
			go addCollectionByContract(wg, req.AccountAddress, "erc1155", api, req.ContractAddress)
			wg.Wait()
		} else {
			if firstCollection.UpdatedAt.Before(time.Now().Add(-time.Duration(global.CONFIG.NFT.CacheTime) * time.Minute)) {
				wg := new(sync.WaitGroup)
				wg.Add(2)
				go addCollectionByContract(wg, req.AccountAddress, "erc721", api, req.ContractAddress)
				go addCollectionByContract(wg, req.AccountAddress, "erc1155", api, req.ContractAddress)
				wg.Wait()
			}
		}
	}
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	db := global.DB.Model(&model.Collection{}).Where("account_address", req.AccountAddress).Where("contract_address", req.ContractAddress).Where("chain", api.Chain)
	err = db.Count(&total).Error
	if err != nil {
		return total, res, err
	}
	err = db.Limit(limit).Offset(offset).Order("own_timestamp desc").Find(&res).Error
	if err != nil {
		return total, res, err
	}

	return
}

// AddCollection
// @description: add collection by ids
// @param: ids []string, address string
// @return: err error
func AddCollection(address string, req request.AddCollectionReq) (err error) {
	tx := global.DB.Begin()
	// 添加NFT
	err = tx.Model(&model.Collection{}).
		Where("chain = ? AND contract_address = ? AND status = 0", req.Chain, req.ContractAddress).
		Updates(map[string]interface{}{"status": 2}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	// 隐藏NFT
	for _, id := range req.HideIDs {
		if id == "" {
			continue
		}
		err = tx.Model(&model.Collection{}).Where("id", id).Updates(map[string]interface{}{"status": 1}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 显示NFT
	for _, id := range req.ShowIDs {
		if id == "" {
			continue
		}
		err = tx.Model(&model.Collection{}).Where("id", id).Updates(map[string]interface{}{"status": 2}).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	// 添加合约到用户
	var idContract string
	errNFTContract := tx.Model(&model.Contract{}).Select("id").Where("chain", req.Chain).
		Where("contract_address", req.ContractAddress).First(&idContract).Error
	if errNFTContract != nil {
		tx.Rollback()
		return err
	}
	_ = addContractToUser(idContract, address)
	// 更新合约数量
	updateContractCount(address)
	return tx.Commit().Error
}

// updateAllCollection
// @description: update account all collection
// @param: address string
func updateAllCollection(address string, uuidList []string, init bool, refresh bool) {
	var count int64
	db := global.DB.Model(&model.Account{})
	limitTime := time.Now().Add(-time.Duration(global.CONFIG.NFT.CacheTime) * time.Minute)
	if refresh {
		limitTime = time.Now().Add(-time.Duration(1) * time.Minute)
	}
	db.Where("(updated_at < ?) AND total < 100 AND address = ? ", limitTime, address)
	err := db.Count(&count).Error
	if err != nil {
		global.LOG.Error("error getting", zap.Error(err))
		return
	}
	if count == 0 && !init {
		return
	}
	// 更新时间
	err = global.DB.Model(&model.Account{}).
		Where("address", address).Update("updated_at", time.Now()).Error
	if err != nil {
		global.LOG.Error("error update", zap.Error(err))
		return
	}
	var total int
	for _, api := range global.CONFIG.NFT.APIConfig {
		t, _ := addAllCollection(address, api, "")
		total += t
	}
	err = global.DB.Model(&model.Account{}).
		Where("address", address).Update("total", total).Error
	if err != nil {
		global.LOG.Error("error update", zap.Error(err))
		return
	}
	for _, v := range uuidList {
		var contracts string
		if err = global.DB.Model(&model.Contract{}).Select("contract_address").Where("id", v).First(&contracts).Error; err != nil {
			return
		}
		if err = global.DB.Model(&model.Collection{}).Where("account_address", address).
			Where("contract_address", contracts).
			Where("status=1").
			Updates(map[string]interface{}{"status": 2}).Error; err != nil {
			return
		}
	}
	// 更新用户展示数量
	_ = updateContractCount(address)
}

// addCollection
// @description: add or update collection to account
// @param: address string, erc_type string, api config.APIConfig
// @return: err error
func addCollectionByContract(wg *sync.WaitGroup, address string, erc_type string, api config.APIConfig, contract string) (err error) {
	defer wg.Done()
	var user model.Account
	contractList := make(map[common.Address]struct{})
	contractList[common.HexToAddress(contract)] = struct{}{}
	if err = global.DB.Model(&model.Account{}).Select("contract_ids").Where("address", address).First(&user).Error; err != nil {
		return err
	}
	for _, v := range user.ContractIDs {
		var contracts string
		if err = global.DB.Model(&model.Contract{}).Select("contract_address").Where("id", v).First(&contracts).Error; err != nil {
			return err
		}
		contractList[common.HexToAddress(contracts)] = struct{}{}
	}
	assetsUrl := fmt.Sprintf("https://%s.nftscan.com/api/v2/account/own/%s?limit=300&erc_type=%s&contract_address=%s", api.APIPreHost, address, erc_type, contract)
	var cursor string
	var nft []model.Collection
	var errFlag bool
	client := req.C().SetTimeout(120*time.Second).
		SetCommonRetryCount(1).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("Referer", "https://docs.nftscan.com/")
	if erc_type == "erc721" {
		client.SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKey)
	} else {
		client.SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKeyBackup)
	}
	i := 0
	for {
		fmt.Println("i", i)
		i++
		fmt.Println(time.Now())
		reqUrl := assetsUrl + fmt.Sprintf("&cursor=%s", cursor)
		fmt.Println(reqUrl)

		req, errReq := client.R().Get(reqUrl)
		if errReq != nil {
			fmt.Println(errReq)
			errFlag = true
			break
		}
		fmt.Println(time.Now())
		res := req.String()

		if gjson.Get(res, "data.total").Uint() == 0 {
			fmt.Println("No data")
			break
		}

		if gjson.Get(res, "code").String() != "200" {
			errFlag = true
			break
		}
		var nftScan []model.NFTScanOwn
		if errParse := json.Unmarshal([]byte(gjson.Get(res, "data.content").String()), &nftScan); errParse != nil {
			fmt.Println(errParse)
			errFlag = true
			break
		}

		for _, v := range nftScan {
			if _, ok := contractList[common.HexToAddress(v.ContractAddress)]; !ok {
				continue
			}
			nft = append(nft, model.Collection{Chain: api.Chain, AccountAddress: address, NFTScanOwn: v})
		}
		cursor = gjson.Get(res, "data.next").String()
		if cursor == "" {
			break
		}
	}

	if len(nft) == 0 {
		return nil
	}
	// 保存数据
	if err = global.DB.Model(&model.Collection{}).Omit("status").Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}, {Name: "account_address"}, {Name: "contract_address"}, {Name: "token_id"}},
		UpdateAll: true,
	}).Create(&nft).Error; err != nil {
		return err
	}
	_ = errFlag
	//if !errFlag {
	//	go filtrateNFT(address, &nft, api)
	//}
	// get item details
	temp := make(map[string]struct{})
	for _, v := range nft {
		if _, ok := temp[v.ContractAddress]; !ok {
			temp[v.ContractAddress] = struct{}{}
		}
	}
	if len(temp) == 0 {
		return
	}
	go ItemFiltrateAndDown(temp, api)

	return nil
}

// addCollection
// @description: add or update collection to account
// @param: address string, erc_type string, api config.APIConfig
// @return: err error
func addAllCollection(address string, api config.APIConfig, contract string) (total int, err error) {
	defer func() {
		if e := recover(); e != nil {
			_ = global.DB.Model(&model.Account{}).Where("address", address).Update("special", true).Error
			return
		}
	}()
	var user model.Account
	contractList := make(map[common.Address]struct{})
	contractList[common.HexToAddress(contract)] = struct{}{}

	if err = global.DB.Model(&model.Account{}).Select("contract_ids").Where("address", address).First(&user).Error; err != nil {
		return total, err
	}
	// TODO 默认合约
	var contractDefault []string
	err = global.DB.Model(&model.ContractDefault{}).Raw("SELECT contract_id FROM contract_default").Scan(&contractDefault).Error
	dealList := slice.DiffSlice[string](contractDefault, user.ContractIDs)

	for _, v := range append(user.ContractIDs, dealList...) {
		var contracts string
		if err = global.DB.Model(&model.Contract{}).Select("contract_address").Where("id", v).First(&contracts).Error; err != nil {
			return total, err
		}
		contractList[common.HexToAddress(contracts)] = struct{}{}
	}
	assetsUrl := fmt.Sprintf("https://%s.nftscan.com/api/v2/account/own/all/%s?show_attribute=false", api.APIPreHost, address)
	//var cursor string
	var nft []model.Collection
	var tryTimes uint
	var errFlag bool
	client := req.C().SetTimeout(500*time.Second).
		SetCommonRetryCount(1).
		SetUserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/111.0.0.0 Safari/537.36").
		SetCommonHeader("Referer", "https://docs.nftscan.com/").
		SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKey)

	//reqUrl := assetsUrl + fmt.Sprintf("&cursor=%s", cursor)
	reqUrl := assetsUrl
	fmt.Println(reqUrl)

	req, errReq := client.R().Get(assetsUrl)
	if errReq != nil {
		tryTimes += 1
		errFlag = true
	}
	req.ErrorResult()
	res := req.String()

	if gjson.Get(res, "code").String() != "200" {
		tryTimes += 1
		errFlag = true
	}
	var nftScan []model.NFTScanOwn
	arr := gjson.Get(res, "data.#.assets").Array()
	for i, _ := range arr {
		var v []model.NFTScanOwn
		if errParse := json.Unmarshal([]byte(arr[i].String()), &v); errParse != nil {
			tryTimes += 1
			errFlag = true
		}
		nftScan = append(nftScan, v...)
	}
	total = len(nftScan)
	for _, v := range nftScan {
		if _, ok := contractList[common.HexToAddress(v.ContractAddress)]; !ok {
			continue
		}
		nft = append(nft, model.Collection{Chain: api.Chain, AccountAddress: address, NFTScanOwn: v, Status: 2})
	}

	if len(nft) == 0 {
		return total, nil
	}
	// 保存数据
	if err = global.DB.Model(&model.Collection{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}, {Name: "account_address"}, {Name: "contract_address"}, {Name: "token_id"}},
		DoNothing: true,
	}).Create(&nft).Error; err != nil {
		return total, err
	}
	_ = errFlag
	// 删除非本人NFT
	//if !errFlag {
	//	go filtrateNFT(address, &nft, api)
	//}
	// get item details
	temp := make(map[string]struct{})
	for _, v := range nft {
		if _, ok := temp[v.ContractAddress]; !ok {
			temp[v.ContractAddress] = struct{}{}
		}
	}
	if len(temp) == 0 {
		return
	}
	go ItemFiltrateAndDown(temp, api)

	return total, nil
}

// filtrateNFT
// @description: if collection not this account，remove collection in this account
// @param: address string, erc_type string, api config.APIConfig
// @return: err error
func filtrateNFT(address string, res *[]model.Collection, api config.APIConfig) (err error) {
	var nftSlice []string
	var ids []string
	for _, v := range *res {
		ids = append(ids, v.ID)
	}
	global.DB.Model(&model.Collection{}).Select("id").Where("chain", api.Chain).
		Where("account_address", address).
		Find(&nftSlice)
	dealList := slice.DiffSlice[string](nftSlice, ids)
	if len(dealList) == 0 {
		return nil
	}
	err = global.DB.Model(&model.Collection{}).Where("id", dealList).Delete(&model.Collection{}).Error
	return err
}

func RefreshUserData(address string) (err error) {
	if address == common.HexToAddress("0").String() {
		return nil
	}
	db := global.DB
	var user model.Account
	var dealList []string
	errFirst := db.Model(&model.Account{}).Where("address = ?", address).First(&user).Error
	if errFirst != nil && errFirst != gorm.ErrRecordNotFound {
		return errFirst
	}

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
		updateAllCollection(address, dealList, false, true) // update all collection
	}

	if len(user.ContractIDs) != len(user.Counts) {
		updateContractCount(address)
	}
	return nil
}
