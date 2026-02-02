package domain

import "time"

type NotificationConfig struct {
	RateSpikeThreshold float64
	SlowThresold       time.Duration
}

func DefaultNotificationConfig() NotificationConfig {
	return NotificationConfig{
		RateSpikeThreshold: 0.1, // 10%
		SlowThresold:       100 * time.Millisecond,
	}
}
