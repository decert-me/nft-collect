package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model"
	"nft-collect/internal/app/model/request"
)

// SaveCardInfo 保存Zcloak证书
func SaveCardInfo(c *gin.Context, r request.SaveCardInfoRequest) (err error) {
	// 校验Key
	if c.GetHeader("x-api-key") != global.CONFIG.System.APIKey {
		global.LOG.Error("非法请求", zap.String("x-api-key", c.GetHeader("x-api-key")))
		return errors.New("非法请求")
	}
	// 查询合约信息
	//var contract model.Contract
	//err = global.DB.Model(&model.Contract{}).Where("contract_address = ?", r.ContractAddress).First(&contract).Error
	//if err != nil {
	//	return
	//}
	// 查询NFT信息

	// 写入NFT数据
	collection := model.Collection{
		Chain:          r.Chain,
		AccountAddress: r.AccountAddress,
		Status:         2, // 显示
		IsZcloak:       true,
		NFTScanOwn: model.NFTScanOwn{
			ContractAddress: r.ContractAddress,
			TokenID:         r.TokenID,
			Owner:           r.AccountAddress,
			Name:            r.Name,
			ErcType:         r.ErcType,
			ImageURI:        r.ImageURI,
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
