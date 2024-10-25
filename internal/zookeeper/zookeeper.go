package zookeeper

import (
	"time"

	"github.com/go-zookeeper/zk"
	"github.com/starryrbs/watchdog/internal/checker"
	"go.uber.org/zap"
)

// Config ZookeeperConfig 配置 Zookeeper 连接
type Config struct {
	Hosts   []string
	Timeout time.Duration
	logger  *zap.Logger
	conn    *zk.Conn
}

type Checker struct {
	logger *zap.Logger
	conn   *zk.Conn
	config *Config
}

func (z *Checker) Name() string {
	return "zookeeper"
}

func NewChecker(config *Config, logger *zap.Logger) checker.Checker {
	return &Checker{
		logger: logger,
		config: config,
	}
}

func (z *Checker) InitConnection() error {
	conn, _, err := zk.Connect(z.config.Hosts, z.config.Timeout)
	if err != nil {
		return err
	}
	z.conn = conn
	return nil

}

func (z *Checker) CheckAvailability() bool {
	_, _, err := z.conn.Get("/")
	if err != nil {
		z.logger.Error("zookeeper connection failed", zap.Error(err))
		return false
	}
	return true
}
