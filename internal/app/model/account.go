package model

import (
	"github.com/lib/pq"
	"nft-collect/internal/app/global"
)

type Account struct {
	global.MODEL
	Address     string         `gorm:"type:char(42);index:account_address,unique;not null;"  json:"address" form:"address"`
	ContractIDs pq.StringArray `gorm:"type:uuid[]" json:"contract_ids" form:"contract_ids"`  // 合约ID
	Counts      pq.Int64Array  `gorm:"type:integer[]" json:"counts" form:"counts"`           // 数量
	CountsShow  pq.Int64Array  `gorm:"type:integer[]" json:"counts_show" form:"counts_show"` // 显示数量
	//CountsDefault pq.Int64Array  `gorm:"type:integer[]" json:"counts_default" form:"counts_default"` // 默认合约数量
	Total int `gorm:"column:total;default:0"` // NFT总数
}
