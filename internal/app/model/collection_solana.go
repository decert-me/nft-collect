package model

import (
	"gorm.io/datatypes"
	"nft-collect/internal/app/global"
)

type CollectionSolana struct {
	global.MODEL
	Status      uint8 `gorm:"default:0;" json:"status" form:"status"`     // 显示状态(0:初始状态 1:隐藏 2:显示)
	ClaimStatus uint8 `gorm:"claim_status;default:1" json:"claim_status"` // 0 未领取 1 NFT 2 zcloak 3 两者
	NFTScanSolana
}

type NFTScanSolana struct {
	Collection          string         `gorm:"column:collection" json:"collection" form:"collection"`          // Collection 地址
	TokenAddress        string         `gorm:"column:token_address" json:"token_address" form:"token_address"` // Token 地址
	Minter              string         `json:"minter" form:"minter"`
	Owner               string         `json:"owner" form:"owner"`
	MintTimestamp       int64          `json:"mint_timestamp"`
	MintTransactionHash string         `json:"mint_transaction_hash"`
	TokenURI            string         `json:"token_uri"`
	MetadataJSON        datatypes.JSON `json:"metadata_json"`
	Name                string         `json:"name"`
	ContentType         string         `json:"content_type"`
	ContentURI          string         `json:"content_uri"`
	ImageURI            string         `json:"image_uri"`
	ExternalLink        string         `json:"external_link"`
	TokenID             string         `json:"token_id" form:"token_id"`
}
