package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/service"
	"nft-collect/internal/app/utils"
	"strings"
)

func GetCollectionByContract(c *gin.Context) {
	var req request.GetCollectionReq
	_ = c.ShouldBindQuery(&req)
	req.ContractAddress = strings.ToLower(c.Param("address"))
	req.AccountAddress = c.GetString("address")
	if total, list, err := service.GetCollectionByContract(req); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, "Success", c)
	}
}
func GetCollection(c *gin.Context) {
	var req request.GetCollectionReq
	_ = c.ShouldBindQuery(&req)
	address := c.Param("address")

	account := c.GetString("address")
	// 检验字段
	if err := utils.Verify(req.PageInfo, utils.PageSizeLimitVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if utils.IsValidAddress(address) {
		req.AccountAddress = strings.ToLower(address)
		if total, totalPublic, totalHidden, list, err := service.GetCollection(req, account); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error", c)
		} else {
			response.OkWithDetailed(response.GetCollectionRes{
				List:        list,
				Total:       total,
				TotalPublic: totalPublic,
				TotalHidden: totalHidden,
				Page:        req.Page,
				PageSize:    req.PageSize,
			}, "Success", c)
		}
	} else if utils.IsValidSolanaAddress(address) {
		req.AccountAddress = address
		if total, totalPublic, totalHidden, list, err := service.GetSolanaCollection(req, account); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error", c)
		} else {
			response.OkWithDetailed(response.GetCollectionRes{
				List:        list,
				Total:       total,
				TotalPublic: totalPublic,
				TotalHidden: totalHidden,
				Page:        req.Page,
				PageSize:    req.PageSize,
			}, "Success", c)
		}
	} else {
		response.FailWithMessage("地址错误", c)
	}
}

func GetContract(c *gin.Context) {
	address := c.Param("address")
	account := c.GetString("address")
	if utils.IsValidAddress(address) {
		address = strings.ToLower(address)
		if list, err := service.GetContract(address, account); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error", c)
		} else {
			response.OkWithDetailed(list, "Success", c)
		}
	} else if utils.IsValidSolanaAddress(address) {
		//response.OkWithDetailed(nil, "Success", c)
		//return
		if list, err := service.GetSolanaContract(address, account); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error", c)
		} else {
			response.OkWithDetailed(list, "Success", c)
		}
	} else {
		response.FailWithMessage("地址错误", c)
	}
}

func AddCollection(c *gin.Context) {
	var req request.AddCollectionReq
	_ = c.ShouldBindJSON(&req)
	address := c.GetString("address")
	req.Chain = global.ChainName[req.ChainID]
	if address == "" || req.Chain == "" {
		response.FailWithMessage("Error", c)
		return
	}

	if err := service.AddCollection(address, req); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithMessage("Success", c)
	}
}

func UpdatedCollection(c *gin.Context) {
	var req request.UpdatedCollectionReq
	_ = c.ShouldBindJSON(&req)
	address := c.GetString("address")
	req.ID = c.Param("id")
	if len(req.ID) == 0 || address == "" {
		response.FailWithMessage("Error", c)
		return
	}
	if err := service.UpdatedCollection(req, address); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithMessage("Success", c)
	}
}

func RefreshUserData(c *gin.Context) {
	var req request.RefreshUserDataReq
	_ = c.ShouldBindJSON(&req)
	address := req.Address

	if utils.IsValidAddress(address) {
		req.Address = strings.ToLower(req.Address)
		if err := service.RefreshUserData(req.Address); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error", c)
		} else {
			response.OkWithMessage("Success", c)
		}
	} else if utils.IsValidSolanaAddress(address) {
		service.RefreshUserDataSolana()
		response.OkWithMessage("Success", c)
	} else {
		response.FailWithMessage("地址错误", c)
	}

}
