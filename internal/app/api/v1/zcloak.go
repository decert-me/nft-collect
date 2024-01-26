package v1

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nft-collect/internal/app/global"
	"nft-collect/internal/app/model/request"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/service"
	"nft-collect/internal/app/utils"
)

// SaveCardInfo 保存Zcloak证书
func SaveCardInfo(c *gin.Context) {
	var req request.SaveCardInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage("ParameterError", c)
		return
	}
	if utils.IsValidAddress(req.AccountAddress) {
		if err := service.SaveCardInfo(c, req); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error: "+err.Error(), c)
		} else {
			response.OkWithMessage("Success", c)
		}
	} else {
		if err := service.SaveSolanaCardInfo(c, req); err != nil {
			global.LOG.Error("Error!", zap.Error(err))
			response.FailWithMessage("Error: "+err.Error(), c)
		} else {
			response.OkWithMessage("Success", c)
		}
	}

}
