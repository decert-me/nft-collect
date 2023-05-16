package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/service"
)

func GetDefaultContract(c *gin.Context) {
	if list, err := service.GetDefaultContract(); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithDetailed(list, "Success", c)
	}
}

func AddDefaultContract(c *gin.Context) {
	var req request.AddDefaultContractReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("ParameterError", c)
		return
	}
	if err := service.AddDefaultContract(req); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithMessage("Success", c)
	}
}

func DelDefaultContract(c *gin.Context) {
	var req request.DelDefaultContractReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("ParameterError", c)
		return
	}
	if err := service.DelDefaultContract(req); err != nil {
		global.LOG.Error("Error!", zap.Error(err))
		response.FailWithMessage("Error", c)
	} else {
		response.OkWithMessage("Success", c)
	}
}
