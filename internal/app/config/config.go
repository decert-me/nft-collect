package config

type Server struct {
	Zap    Zap    `mapstructure:"zap" json:"zap" yaml:"zap"`
	System System `mapstructure:"system" json:"system" yaml:"system"`
	// gorm
	Pgsql Pgsql `mapstructure:"pgsql" json:"pgsql" yaml:"pgsql"`
	NFT   NFT   `mapstructure:"nft" json:"nft" yaml:"nft"`
	JWT   JWT   `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
}
