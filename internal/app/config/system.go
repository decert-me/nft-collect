package config

type System struct {
	Env    string `mapstructure:"env" json:"env" yaml:"env"`             // 环境值
	Addr   int    `mapstructure:"addr" json:"addr" yaml:"addr"`          // 端口值
	APIKey string `mapstructure:"api-key" json:"api-key" yaml:"api-key"` // API Key
}
