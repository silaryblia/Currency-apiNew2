package config

import (
	"fmt"
	"os"
	"strings"
)

func Load() (Config, error) {
	cfg := DefaultConfig()

	if v := os.Getenv(EnvUsePostgres); v != "" {
		cfg.UsePostgres = v == "true"
	}
	if v := os.Getenv(EnvPostgresDSN); v != "" {
		cfg.PostgresDSN = v
	}
	if v := os.Getenv(EnvHTTPPort); v != "" {
		cfg.HTTPPort = v
	}
	if v := os.Getenv(EnvGRPCPort); v != "" {
		cfg.GRPCPort = v
	}
	if v := os.Getenv(EnvLogMode); v != "" {
		cfg.LogMode = strings.ToLower(v)
	}

	switch cfg.LogMode {
	case "dev", "prod":
		// ok
	default:
		return Config{}, fmt.Errorf("invalid LOG_MODE: %s (allowed: dev, prod)", cfg.LogMode)
	}

	if v := os.Getenv("CBR_URL"); v != "" {
		cfg.CBR.URL = v
	}
	if v := os.Getenv("CBR_REQUIRED_CODES"); v != "" {
		cfg.CBR.RequiredCodes = strings.Split(v, ",")
	}

	if cfg.GRPCPort == "" {
		return Config{}, fmt.Errorf("GRPC_PORT must not be empty")
	}

	if cfg.UsePostgres && cfg.PostgresDSN == "" {
		return Config{}, fmt.Errorf("POSTGRES_DSN is required when USE_POSTGRES=true")
	}

	if cfg.CBR.URL == "" {
		return Config{}, fmt.Errorf("CBR_URL must not be empty")
	}

	if len(cfg.CBR.RequiredCodes) == 0 {
		return Config{}, fmt.Errorf("CBR_REQUIRED_CODES must not be empty")
	}

	return *cfg, nil
}
