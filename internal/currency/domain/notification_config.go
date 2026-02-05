package domain

import "time"

type NotificationConfig struct {
	RateSpikeThreshold float64       `env:"RATE_SPIKE_THRESHOLD"`
	SlowThreshold      time.Duration `env:"SLOW_THRESHOLD"`
}

func DefaultNotificationConfig() NotificationConfig {
	return NotificationConfig{
		RateSpikeThreshold: 0.1, // 10%
		SlowThreshold:      100 * time.Millisecond,
	}
}
