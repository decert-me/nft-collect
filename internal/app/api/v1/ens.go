package v1

import (
	"github.com/gin-gonic/gin"
	"nft-collect/internal/app/model/response"
	"nft-collect/internal/app/service"
)

func GetEnsRecords(c *gin.Context) {
	if list, err := service.GetEnsRecords(c, c.Param("q")); err != nil {
	} else {
		response.OkWithData(list, c)
	}
}
