package request

import "nft-collect/internal/app/model"

type AddContractReq struct {
	ContractAddress string `json:"contract_address" form:"contract_address"` // 合约地址
}

type AddCollectionReq struct {
	Flag int64    `json:"flag"`
	IDs  []string `json:"ids" form:"ids"`
}

type UpdatedCollectionReq struct {
	ID     string `json:"id" form:"id"`
	Status uint8  `gorm:"default:1;" json:"status" form:"status"` // 显示状态(1:隐藏 2:显示)
}

type GetCollectionReq struct {
	PageInfo
	model.Collection
	ContractID string `form:"contract_id"`
	ChainID    uint   `form:"chain_id"`
	Search     string `form:"search"`
	Sort       string `form:"sort"`
}

type AddressReq struct {
	Address string `json:"address" form:"address"`
}

type RefreshUserDataReq struct {
	Address string `json:"address" form:"address"`
}
