package request

import "gorm.io/datatypes"

type SaveCardInfoRequest struct {
	Chain           string         `json:"chain" form:"chain" binding:"required"`
	AccountAddress  string         `json:"account_address" form:"account_address" binding:"required"`
	ContractAddress string         `json:"contract_address" form:"contract_address" binding:"required"`
	TokenID         string         `json:"token_id" form:"token_id" binding:"required"`
	ImageURI        string         `json:"image_uri" form:"image_uri" binding:"required"`
	ErcType         string         `json:"erc_type" form:"erc_type" binding:"required"`
	Name            string         `json:"name" form:"name" binding:"required"`
	DidAddress      string         `json:"did_address" form:"did_address" binding:"required"`
	MetadataJson    datatypes.JSON `json:"metadata_json" form:"metadata_json"`
	UUID            string         `json:"uuid" form:"uuid" binding:"required"`
}
