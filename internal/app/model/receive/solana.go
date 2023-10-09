package receive

type SolanaCollection struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	Data struct {
		Total   int         `json:"total"`
		Next    interface{} `json:"next"`
		Content []struct {
			BlockNumber                int         `json:"block_number"`
			InteractProgram            string      `json:"interact_program"`
			Collection                 string      `json:"collection"`
			TokenAddress               string      `json:"token_address"`
			Minter                     string      `json:"minter"`
			Owner                      string      `json:"owner"`
			MintTimestamp              int64       `json:"mint_timestamp"`
			MintTransactionHash        string      `json:"mint_transaction_hash"`
			MintPrice                  int         `json:"mint_price"`
			TokenURI                   string      `json:"token_uri"`
			MetadataJSON               string      `json:"metadata_json"`
			Name                       string      `json:"name"`
			ContentType                string      `json:"content_type"`
			ContentURI                 string      `json:"content_uri"`
			ImageURI                   string      `json:"image_uri"`
			ExternalLink               string      `json:"external_link"`
			LatestTradePrice           interface{} `json:"latest_trade_price"`
			LatestTradeTimestamp       interface{} `json:"latest_trade_timestamp"`
			LatestTradeTransactionHash interface{} `json:"latest_trade_transaction_hash"`
		} `json:"content"`
	} `json:"data"`
}
