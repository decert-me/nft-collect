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
	req.AccountAddress = strings.ToLower(c.Param("address"))
	account := c.GetString("address")
	// 检验字段
	if err := utils.Verify(req.PageInfo, utils.PageSizeLimitVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
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
}

func GetContract(c *gin.Context) {
	address := strings.ToLower(c.Param("address"))
	account := c.GetString("address")
	if !utils.IsValidAddress(address) {
		response.FailWithMessage("Param Error", c)
		return
	}
	if list, err := service.GetContract(address, account); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithDetailed(list, "Success", c)
	}
}

func AddCollection(c *gin.Context) {
	var req request.AddCollectionReq
	_ = c.ShouldBindJSON(&req)
	address := c.GetString("address")
	if len(req.IDs) == 0 || address == "" {
		response.FailWithMessage("Error", c)
		return
	}
	if err := service.AddCollection(req.IDs, address, req.Flag); err != nil {
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
	if err := service.RefreshUserData(req.Address); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithMessage("Success", c)
	}
}
