package config

import "strings"

type DBLabel string

type MysqlNodeConfig struct {
	DSN                    string `yaml:"dsn" mapstructure:"dsn"`
	MaxOpenConns           int    `yaml:"max_open_conns" mapstructure:"max_open_conns"`
	MaxIdleConns           int    `yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	ConnMaxLifetimeSeconds int    `yaml:"conn_max_lifetime_seconds" mapstructure:"conn_max_lifetime_seconds"`
	ConnMaxIdleTimeSeconds int    `yaml:"conn_max_idle_time_seconds" mapstructure:"conn_max_idle_time_seconds"`
}

type MySQLClusterConfig struct {
	DBLabel  DBLabel           `yaml:"db_label" mapstructure:"db_label"`
	Master   MysqlNodeConfig   `yaml:"master" mapstructure:"master"`
	Replicas []MysqlNodeConfig `yaml:"replicas" mapstructure:"replicas"`
}

type DBParams struct {
	CIDs []string `yaml:"cids" mapstructure:"cids"`
	ENV  string   `yaml:"env" mapstructure:"env"`
}

type MySQLConfig struct {
	DBParams           DBParams `yaml:"db_params" mapstructure:"db_params"`
	EagerLoad          bool     `yaml:"eager_load" mapstructure:"eager_load"`
	MySQLClusterConfig `mapstructure:",squash"`
}

func (m *MySQLConfig) GetMysqlClusterConfig() []MySQLClusterConfig {
	env := m.DBParams.ENV
	buildCluster := func(cid string) MySQLClusterConfig {
		mysqlCluster := m.MySQLClusterConfig
		mysqlCluster.DBLabel = DBLabel(formatDBInfoStr(string(mysqlCluster.DBLabel), cid, env))
		mysqlCluster.Master.DSN = formatDBInfoStr(mysqlCluster.Master.DSN, cid, env)
		for idx := range mysqlCluster.Replicas {
			mysqlCluster.Replicas[idx].DSN = formatDBInfoStr(mysqlCluster.Replicas[idx].DSN, cid, env)
		}
		return mysqlCluster
	}

	if len(m.DBParams.CIDs) == 0 {
		return []MySQLClusterConfig{buildCluster("")}
	}

	mysqlClusterConfigs := make([]MySQLClusterConfig, 0, len(m.DBParams.CIDs))
	for _, cid := range m.DBParams.CIDs {
		mysqlClusterConfigs = append(mysqlClusterConfigs, buildCluster(cid))
	}
	return mysqlClusterConfigs
}

func formatDBInfoStr(str, cid, env string) string {
	if cid != "" {
		str = strings.ReplaceAll(str, "{cid}", cid)
	}
	if env != "" {
		str = strings.ReplaceAll(str, "{env}", env)
	}
	return str
}
