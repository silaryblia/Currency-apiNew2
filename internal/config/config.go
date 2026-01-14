package config

import "time"

type Config struct {
	UsePostgres bool
	PostgresDSN string
	HTTPPort    string
	GRPCPort    string
	LogMode     string

	CBR CBRConfig `yaml:"cbr"`
}
type HTTPConfig struct {
	Timeout         time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT"`
	MaxIdleConns    int           `yaml:"maxIdleConns" env:"HTTP_MAX_IDLE_CONNS"`
	IdleConnTimeout time.Duration `yaml:"idleConnTimeout" env:"HTTP_IDLE_CONN_TIMEOUT"`
}

type CBRConfig struct {
	URL           string     `yaml:"url" env:"CBR_URL"`
	RequiredCodes []string   `yaml:"requiredCodes" env:"CBR_REQUIRED_CODES" envSeparator:","`
	HTTPConfig    HTTPConfig `yaml:"http"`
}

func DefaultCBRConfig() *CBRConfig {
	return &CBRConfig{
		URL:           "https://www.cbr.ru/scripts/XML_daily.asp",
		RequiredCodes: []string{"USD", "EUR", "AED"},
		HTTPConfig: HTTPConfig{
			Timeout:         30 * time.Second,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
	}
}

func DefaultConfig() *Config {
	return &Config{
		UsePostgres: false,
		PostgresDSN: "",
		HTTPPort:    "8080",
		GRPCPort:    "50051",
		LogMode:     "dev",
		CBR:         *DefaultCBRConfig(),
	}
}
