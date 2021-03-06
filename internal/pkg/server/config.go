package server

import (
	"github.com/gin-gonic/gin"
)

type Config struct {
	Mode            string
	Middlewares     []string
	Healthz         bool
	EnableProfiling bool
	EnableMetrics   bool
}

func NewConfig() *Config {
	return &Config{
		Healthz:         true,
		Mode:            gin.ReleaseMode,
		Middlewares:     []string{},
		EnableProfiling: true,
		EnableMetrics:   true,
	}
}

// CompletedConfig is the completed configuration for GenericAPIServer.
type CompletedConfig struct {
	*Config
}

func (c *Config) Complete() CompletedConfig {
	return CompletedConfig{c}
}

func (c CompletedConfig) New() (*GenericAPIServer, error) {
	s := &GenericAPIServer{
		mode:            c.Mode,
		healthz:         c.Healthz,
		enableMetrics:   c.EnableMetrics,
		enableProfiling: c.EnableProfiling,
		middlewares:     c.Middlewares,
		Engine:          gin.New(),
	}

	initGenericAPIServer(s)

	return s, nil
}
