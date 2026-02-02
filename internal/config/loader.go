package config

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var localConfigMu sync.RWMutex

// LoadLocalConfig loads config from file and enables hot reload.
func LoadLocalConfig(configPath string) error {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := updateLocalConfig(v); err != nil {
		return err
	}

	v.OnConfigChange(func(event fsnotify.Event) {
		if err := updateLocalConfig(v); err != nil {
			slog.Error("reload local config failed", "path", configPath, "error", err)
			return
		}
		slog.Info("local config reloaded", "path", configPath, "event", event.String())
	})
	v.WatchConfig()

	return nil
}

func updateLocalConfig(v *viper.Viper) error {
	var cfg LocalConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return err
	}
	fmt.Printf("update cfg:%+v\n", cfg)
	localConfigMu.Lock()
	LocalCfg = cfg
	localConfigMu.Unlock()

	return nil
}
