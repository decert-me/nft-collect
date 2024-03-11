package config

type NFT struct {
	EnsRpc      string      `mapstructure:"ens-rpc" json:"ens-rpc" yaml:"ens-rpc"` // ENS 查询 RPC
	ApiKey      string      `mapstructure:"api-key" json:"api-key" yaml:"api-key"`
	CacheTime   int         `mapstructure:"cache-time" json:"cache-time" yaml:"cache-time"` // 缓存时间 分钟
	LogoPath    string      `mapstructure:"logo-path" json:"logo-path" yaml:"logo-path"`
	DefContract []string    `mapstructure:"def-contract" json:"def-contract" yaml:"def-contract"`
	APIConfig   []APIConfig `mapstructure:"api-config" json:"api-config" yaml:"api-config"`
}

type APIConfig struct {
	Chain      string `mapstructure:"chain" json:"chain" yaml:"chain"`
	ChainID    uint   `mapstructure:"chain-id" json:"chain-id" yaml:"chain-id"`
	APIPreHost string `mapstructure:"api-per-host" json:"api-per-host" yaml:"api-per-host"`
	Symbol     string `mapstructure:"symbol" json:"symbol" yaml:"symbol"`
}
