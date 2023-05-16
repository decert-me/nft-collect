package receive

type GetNftPlatformInfo struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		Image           string `json:"image"`
		Website         string `json:"website"`
		Address         string `json:"address"`
		Description     string `json:"description"`
		Banner          string `json:"banner"`
		OpenseaVerify   bool   `json:"opensea_verify"`
		Royalty         string `json:"royalty"`
		AuthFlag        bool   `json:"authFlag"`
		PageView        int    `json:"pageView"`
		Name            string `json:"name"`
		ContractCreator string `json:"contractCreator"`
	} `json:"data"`
}

type GetItemModel struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		ContractAddress         string        `json:"contract_address"`
		Name                    string        `json:"name"`
		Symbol                  string        `json:"symbol"`
		Description             string        `json:"description"`
		Website                 string        `json:"website"`
		Email                   interface{}   `json:"email"`
		Twitter                 string        `json:"twitter"`
		Discord                 string        `json:"discord"`
		Telegram                interface{}   `json:"telegram"`
		Github                  interface{}   `json:"github"`
		Instagram               interface{}   `json:"instagram"`
		Medium                  interface{}   `json:"medium"`
		LogoURL                 string        `json:"logo_url"`
		BannerURL               string        `json:"banner_url"`
		FeaturedURL             string        `json:"featured_url"`
		LargeImageURL           string        `json:"large_image_url"`
		Attributes              []interface{} `json:"attributes"`
		ErcType                 string        `json:"erc_type"`
		DeployBlockNumber       int           `json:"deploy_block_number"`
		Owner                   string        `json:"owner"`
		Verified                bool          `json:"verified"`
		OpenseaVerified         bool          `json:"opensea_verified"`
		Royalty                 int           `json:"royalty"`
		ItemsTotal              int           `json:"items_total"`
		AmountsTotal            int           `json:"amounts_total"`
		OwnersTotal             int           `json:"owners_total"`
		OpenseaFloorPrice       float64       `json:"opensea_floor_price"`
		FloorPrice              float64       `json:"floor_price"`
		CollectionsWithSameName []string      `json:"collections_with_same_name"`
		PriceSymbol             string        `json:"price_symbol"`
	} `json:"data"`
}
