package model

type ZcloakDid struct {
	ID         uint   `gorm:"primarykey"`
	Address    string `gorm:"column:address;type:varchar(100);comment:钱包地址" json:"address" form:"address"`
	DidAddress string `gorm:"column:did_address;type:varchar(100);comment:DID地址" json:"did_address" form:"did_address"`
}
