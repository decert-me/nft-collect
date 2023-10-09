package model

import (
	"nft-collect/internal/app/global"
)

type ContractSolana struct {
	global.MODEL
	Chain               string `gorm:"default:'solana'"  json:"chain" form:"chain"`                        // 区块链的简称
	ContractAddress     string `gorm:"default:''" json:"contract_address" form:"contract_address"`         // 合约地址
	ContractName        string `gorm:"default:'';not null" json:"contract_name" form:"contract_name"`      // 合约名称
	ContractLogo        string `gorm:"default:''" json:"contract_logo" form:"contract_logo"`               // 合约Logo
	ContractBanner      string `gorm:"default:''" json:"contract_banner" form:"contract_banner"`           // 合约Banner
	ContractDescription string `gorm:"default:''" json:"contract_description" form:"contract_description"` // 合约Description
	ContractWebsite     string `gorm:"default:''" json:"contract_website" form:"contract_website"`         // 合约Website
	ContractOwner       string `gorm:"default:''" json:"contract_owner" form:"contract_owner"`             // 合约Owner
	Status              uint8  `gorm:"default:1;" json:"-" form:"-"`                                       // 显示状态(1:未获取 2:已获取)
}
