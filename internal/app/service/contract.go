package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"io"
	"nft-collect/internal/app/config"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/receive"
	"nft-collect/pkg/slice"
	"os"
	"path/filepath"
	"strings"
)

// ItemFiltrateAndDown
// @description: filtrate done and down new connection item details
// @param: contractMap map[string]struct{}, api config.APIConfig
// @return: err error
func ItemFiltrateAndDown(contractMap map[string]struct{}, api config.APIConfig) (err error) {
	var addressExist []string
	var addressList []string
	for key, _ := range contractMap {
		addressList = append(addressList, key)
	}
	if err = global.DB.Model(&model.Contract{}).Select("contract_address").Where("contract_address IN ? AND Status = 2", addressList).Find(&addressExist).Error; err != nil {
		return err
	}
	// 待处理 slice
	dealList := slice.DiffSlice[string](addressList, addressExist)
	if len(dealList) == 0 {
		return nil
	}
	//fmt.Println("待处理", dealList)
	res := downloadItem(dealList, api)
	if len(res) == 0 {
		return
	}
	// 保存数据
	if err = global.DB.Model(&model.Contract{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain"}, {Name: "contract_address"}},
		UpdateAll: true,
	}).Create(&res).Error; err != nil {
		return err
	}
	return nil
}

// downloadItem
// @description: down new connection item details
// @param: contractAddress []string, api config.APIConfig
// @return: res []model.NFTContract
func downloadItem(contractAddress []string, api config.APIConfig) (res []model.Contract) {
	client := req.C().SetCommonHeader("X-API-KEY", global.CONFIG.NFT.ApiKeyPro)
	clientDown := req.C()
	for _, address := range contractAddress {
		url := fmt.Sprintf("https://%s.nftscan.com/api/v2/collections/%s?show_attribute=false", api.APIPreHost, address)
		fmt.Println(url)
		resp, err := client.R().Get(url)
		if err != nil {
			continue
		}
		if !strings.Contains(resp.GetStatus(), "200") {
			continue
		}
		var itemModel receive.GetItemModel
		if errParse := json.Unmarshal(resp.Bytes(), &itemModel); errParse != nil {
			continue
		}
		p := itemModel.Data

		if p.ContractAddress == "" {
			continue
		}
		temp := model.Contract{Chain: api.Chain, ContractAddress: address, ContractName: p.Name,
			ContractLogo: p.LogoURL, ContractBanner: p.BannerURL, ContractDescription: p.Description,
			ContractWebsite: p.Website, ContractOwner: p.Owner, Status: 2,
		}
		if itemModel.Data.LogoURL != "" {
			temp.ContractLogo, err = downloadLogo(clientDown, temp, api)
			if err != nil {
				global.LOG.Error("downloadLogo error: ", zap.Error(err))
			}
		}
		res = append(res, temp)
	}
	return
}

// downloadLogo
// @description: download logo
// @param: client *req.Client, nftContract model.NFTContract, api config.APIConfig
// @return: res string, err error
func downloadLogo(client *req.Client, nftContract model.Contract, api config.APIConfig) (res string, err error) {
	res = global.CONFIG.NFT.LogoPath + "logo.png"
	baseUrl := nftContract.ContractLogo
	fileName := fmt.Sprintf("%s%s", nftContract.ContractAddress, filepath.Ext(nftContract.ContractLogo))
	dirPath := global.CONFIG.NFT.LogoPath + "/" + api.Chain + "/"
	err = os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		global.LOG.Error("Error MkdirAll", zap.Error(err))
		return res, err
	}
	resp, err := client.R().Get(baseUrl)
	if err != nil || !strings.Contains(resp.GetStatus(), "200") {
		return res, err
	}
	defer resp.Body.Close()
	f, err := os.Create(dirPath + fileName)
	if err != nil {
		global.LOG.Error("Error creating", zap.Error(err))
		return res, err
	}
	body, _ := io.ReadAll(resp.Body)
	_, err = io.Copy(f, bytes.NewReader(body))
	if err != nil {
		global.LOG.Error("Error io.Copy", zap.Error(err))
		return res, err
	}
	f.Close()
	return "/" + dirPath + fileName, nil
}
