package db

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Me1onRind/mr_agent/internal/config"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"gorm.io/gorm"
)

var ErrDBLabelNotFound = errors.New("db label not found")

// registry manages named read/write database clients.
// It is safe for concurrent reads; configs are lazily initialized on first use.
type registry struct {
	mu      sync.RWMutex
	entries map[config.DBLabel]*registryEntry
}

var defaultRegistry = newRegistry()

func newRegistry() *registry {
	return &registry{
		entries: make(map[config.DBLabel]*registryEntry),
	}
}

type registryEntry struct {
	cfg    config.MySQLClusterConfig
	client *MysqlClient
	once   sync.Once
	err    error
}

func (e *registryEntry) initClient(ctx context.Context, label config.DBLabel) error {
	e.once.Do(func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		client, err := NewMySQLClient(ctx, e.cfg.Master, e.cfg.Replicas...)
		if err != nil {
			e.err = fmt.Errorf("init db %s: %w", label, err)
			return
		}
		e.client = client
	})
	return e.err
}

// InitRegistry builds the global registry from label-config pairs.
func InitRegistry(ctx context.Context, configs []config.MySQLConfig) error {
	log := logger.LoggerFromCtx(ctx)
	registry := newRegistry()
	for _, cfg := range configs {
		clusters := cfg.GetMysqlClusterConfig()
		for _, cluster := range clusters {
			log.Info("init db", slog.String("db_label", string(cluster.DBLabel)), slog.Bool("eager_load", cfg.EagerLoad))
			if err := registry.register(cluster, cluster.DBLabel); err != nil {
				return err
			}
			if cfg.EagerLoad {
				if err := registry.initClientForLabel(ctx, cluster.DBLabel); err != nil {
					_ = registry.closeAll()
					return err
				}
			}
		}
	}
	defaultRegistry = registry
	return nil
}

// GetMasterDB returns the master DB by label.
func GetMasterDB(ctx context.Context, label config.DBLabel) *gorm.DB {
	client := defaultRegistry.get(ctx, label)
	return client.Write(ctx)
}

// GetSlaveDB returns the read DB by label.
func GetSlaveDB(ctx context.Context, label config.DBLabel) *gorm.DB {
	client := defaultRegistry.get(ctx, label)
	return client.Read(ctx)
}

func (r *registry) register(cfg config.MySQLClusterConfig, label config.DBLabel) error {
	if label == "" {
		return errors.New("db label is empty")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.entries[label]; exists {
		return fmt.Errorf("db label already exists: %s", label)
	}
	r.entries[label] = &registryEntry{cfg: cfg}
	return nil
}

func (r *registry) get(ctx context.Context, label config.DBLabel) *MysqlClient {
	if label == "" {
		panic(errors.New("db label is empty"))
	}
	r.mu.RLock()
	entry, ok := r.entries[label]
	r.mu.RUnlock()
	if !ok {
		panic(fmt.Errorf("%w: %s", ErrDBLabelNotFound, label))
	}
	if err := entry.initClient(ctx, label); err != nil {
		panic(err)
	}
	if entry.client == nil {
		panic(fmt.Errorf("db client is nil: %s", label))
	}
	return entry.client
}

func (r *registry) initClientForLabel(ctx context.Context, label config.DBLabel) error {
	if label == "" {
		return errors.New("db label is empty")
	}
	r.mu.RLock()
	entry, ok := r.entries[label]
	r.mu.RUnlock()
	if !ok {
		return fmt.Errorf("%w: %s", ErrDBLabelNotFound, label)
	}
	return entry.initClient(ctx, label)
}

func (r *registry) closeAll() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for label, entry := range r.entries {
		if entry.client == nil {
			continue
		}
		if err := entry.client.Close(); err != nil {
			return fmt.Errorf("close db %s: %w", label, err)
		}
	}
	return nil
}
