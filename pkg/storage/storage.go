package storage

import "time"

type (
	Cache interface {
		Set(key string, value string, ttl time.Duration)
		Get(key string) (string, bool)
		Delete(key string)
	}
)
