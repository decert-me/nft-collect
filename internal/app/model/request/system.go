package request

type AddDefaultContractReq struct {
	Chain           string `json:"chain" form:"chain"  binding:"required"`
	ContractAddress string `json:"contract_address" form:"contract_address"  binding:"required"` // 合约地址
}
type DelDefaultContractReq struct {
	ID string `json:"id" from:"id" binding:"required"`
}
