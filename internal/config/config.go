package config

type LocalConfig struct {
	MysqlConfigs []MySQLConfig `yaml:"mysql_configs" mapstructure:"mysql_configs"`
}

// LocalConfig holds the in-memory runtime configuration.
var LocalCfg LocalConfig
