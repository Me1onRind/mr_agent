package db

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/Me1onRind/mr_agent/internal/config"
	"github.com/Me1onRind/mr_agent/internal/infrastructure/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlClient provides manual read/write routing for GORM.
type MysqlClient struct {
	master   *gorm.DB
	replicas []*gorm.DB
	rr       uint64
}

func NewMySQLClient(ctx context.Context, masterCfg config.MysqlNodeConfig, replicaCfgs ...config.MysqlNodeConfig) (*MysqlClient, error) {
	master, err := newMysqlNode(ctx, masterCfg)
	if err != nil {
		return nil, err
	}

	replicas := make([]*gorm.DB, 0, len(replicaCfgs))
	for _, replicaCfg := range replicaCfgs {
		replica, err := newMysqlNode(ctx, replicaCfg)
		if err != nil {
			return nil, err
		}
		replicas = append(replicas, replica)
	}

	return &MysqlClient{
		master:   master,
		replicas: replicas,
	}, nil
}

func (c *MysqlClient) Write(ctx context.Context) *gorm.DB {
	return c.master.WithContext(ctx)
}

func (c *MysqlClient) Read(ctx context.Context) *gorm.DB {
	if len(c.replicas) == 0 {
		return c.master.WithContext(ctx)
	}
	idx := atomic.AddUint64(&c.rr, 1)
	return c.replicas[int(idx)%len(c.replicas)].WithContext(ctx)
}

func (c *MysqlClient) Close() error {
	if err := closeGormDB(c.master); err != nil {
		return err
	}
	for _, replica := range c.replicas {
		if err := closeGormDB(replica); err != nil {
			return err
		}
	}
	return nil
}

func newMysqlNode(_ context.Context, cfg config.MysqlNodeConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger: logger.NewGormLogger(logger.GormLoggerOptions{
			SlowThreshold:             200 * time.Millisecond,
			IgnoreRecordNotFoundError: true,
		}),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetimeSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeSeconds) * time.Second)
	}
	if cfg.ConnMaxIdleTimeSeconds > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.ConnMaxIdleTimeSeconds) * time.Second)
	}
	return db, nil
}

func closeGormDB(db *gorm.DB) error {
	if db == nil {
		return nil
	}
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
