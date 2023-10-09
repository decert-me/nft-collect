package response

import (
	"nft-collect/internal/app/model"
)

type GetSolanaContractRes struct {
	model.ContractSolana
	Count int64 `gorm:"-" json:"count" form:"count"`
}

type GetContractRes struct {
	model.Contract
	Count int64 `gorm:"-" json:"count" form:"count"`
}

type GetCollectionRes struct {
	List        interface{} `json:"list"`
	Total       int64       `json:"total"`
	TotalPublic int64       `json:"total_public"`
	TotalHidden int64       `json:"total_hidden"`
	Page        int         `json:"page"`
	PageSize    int         `json:"pageSize"`
}

type GetCollection struct {
	model.Collection
	ContractLogo string `json:"contract_logo" form:"contract_logo"` // 合约Logo
}
