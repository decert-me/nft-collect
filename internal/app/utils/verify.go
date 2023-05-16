package utils

var (
	IdVerify            = Rules{"ID": {NotEmpty()}}
	PageInfoVerify      = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty(), Le("30")}}
	PageSizeLimitVerify = Rules{"PageSize": {Le("30")}}
	// 用户
	AddContractVerify   = Rules{"ContractAddress": {NotEmpty()}}
	GetCollectionVerify = Rules{"AccountAddress": {NotEmpty()}, "PageSize": {Le("30")}}
)
