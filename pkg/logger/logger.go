package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	ModeDev  = "dev"
	ModeProd = "prod"
)

func New(mode string) (*zap.Logger, error) {
	switch mode {
	case ModeDev:
		return zap.NewDevelopment()
	case ModeProd:
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("unknown logger mode: %s", mode)
	}
}
