package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/starryrbs/watchdog/internal/checker"
	"go.uber.org/zap"
)

// Config 表示 Redis 的配置项
type Config struct {
	Addr     string        `yaml:"addr"`
	Password string        `yaml:"password"`
	DB       int           `yaml:"db"`
	Timeout  time.Duration `yaml:"timeout"`
}

type Checker struct {
	logger *zap.Logger
	client *redis.Client
	config *Config
}

func (c *Checker) Name() string {
	return "redis"
}

func NewChecker(config *Config, logger *zap.Logger) checker.Checker {
	return &Checker{
		logger: logger,
		client: nil,
		config: config,
	}
}

func (c *Checker) InitConnection() error {
	c.client = redis.NewClient(&redis.Options{
		Addr:     c.config.Addr,
		Password: c.config.Password,
		DB:       c.config.DB,
	})
	_, err := c.client.Ping(context.Background()).Result()
	return err
}

func (c *Checker) CheckAvailability() bool {
	if _, err := c.client.Ping(context.Background()).Result(); err != nil {
		c.logger.Error("redis.Client ping failed", zap.Error(err))
		return false
	}
	return true
}
