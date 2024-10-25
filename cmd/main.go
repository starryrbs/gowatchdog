package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/starryrbs/watchdog/internal/checker"
	"github.com/starryrbs/watchdog/internal/redis"
	"github.com/starryrbs/watchdog/internal/zookeeper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis     *redis.Config
	Zookeeper *zookeeper.Config
	Interval  time.Duration
}

func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 YAML 失败: %w", err)
	}
	return &config, nil
}

func RunChecker(checker checker.Checker, logger *zap.Logger, config *Config) {
	clogger := logger.With(zap.String("checker", checker.Name()))
	err := checker.InitConnection()
	if err != nil {
		clogger.Error("c.InitConnection failed", zap.Error(err))
		return
	}

	for {
		success := checker.CheckAvailability()
		if success {
			clogger.Info("is available")
		} else {
			clogger.Error("is not available")
		}
		time.Sleep(config.Interval)

	}
}

func main() {
	logger, _ := zap.NewProduction()

	config, err := LoadConfig("config.yaml")
	if err != nil {
		logger.Error("LoadConfig failed", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}

	if config.Zookeeper != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunChecker(zookeeper.NewChecker(config.Zookeeper, logger), logger, config)
		}()
	}

	if config.Redis != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			RunChecker(redis.NewChecker(config.Redis, logger), logger, config)
		}()
	}

	wg.Wait()

}
