package runnel

import (
	"gopkg.in/redis.v5"
)

// Runnel details
type Runnel struct {
	Redis *redis.Client
}

// NewClient to execute user command
func NewClient() *Runnel {
	runnel := &Runnel{}
	runnel.Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return runnel
}
