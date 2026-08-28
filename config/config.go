package config

// SlsConfig 阿里云日志服务（SLS）后端配置，由 log-sls 模块持有，core 不感知。
type SlsConfig struct {
	Endpoint    string `mapstructure:"endpoint"`
	AccessKey   string `mapstructure:"access_key"`
	SecretKey   string `mapstructure:"secret_key"`
	SecretToken string `mapstructure:"secret_token"`
	Project     string `mapstructure:"project"`
	LogStore    string `mapstructure:"log_store"`
}
