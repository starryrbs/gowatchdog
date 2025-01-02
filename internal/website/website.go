package website

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/starryrbs/watchdog/internal/checker"
	"go.uber.org/zap"
)

type Config struct {
	URLS    []string
	Timeout time.Duration
}

type Checker struct {
	logger *zap.Logger
	config *Config
}

func NewChecker(config *Config, logger *zap.Logger) checker.Checker {
	return &Checker{
		logger: logger,
		config: config,
	}
}

func (c *Checker) InitConnection() error {
	c.logger.Info("website.Checker InitConnection success")
	return nil
}

func (c *Checker) CheckAvailability() bool {
	for _, url := range c.config.URLS {
		err := c.checkSSL(url)
		if err != nil {
			c.logger.Error("website.Checker checkSSL failed", zap.Error(err))
			return false
		}
		c.logger.Info("website.Checker checkSSL success", zap.String("url", url))
	}
	return true
}

func (c *Checker) checkSSL(url string) error {
	client := &http.Client{
		Timeout: c.config.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Validate certificate
			},
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to perform HTTPS request: %v", err)
	}
	defer resp.Body.Close()

	if resp.TLS == nil || len(resp.TLS.PeerCertificates) == 0 {
		return fmt.Errorf("no certificates found")
	}

	// Check expiration date
	now := time.Now()
	for _, cert := range resp.TLS.PeerCertificates {
		if now.After(cert.NotAfter) {
			return fmt.Errorf("certificate expired on %s", cert.NotAfter)
		}
		if now.Before(cert.NotBefore) {
			return fmt.Errorf("certificate not valid until %s", cert.NotBefore)
		}
	}

	return nil
}

func (c *Checker) Name() string {
	return "website"
}
