package runnel

import (
	"fmt"
	"log"
	"time"
)

// Stream output
func (run *Runnel) Stream(key string, timeout int) (string, error) {

	if key == "" {
		return "", fmt.Errorf("Empty key")
	}

	if timeout == 0 {
		timeout = 10
	}

	val, err := run.Redis.BRPop(
		(time.Duration(timeout) * time.Second),
		key).Result()

	if err != nil {
		log.Println(err)
		return "", fmt.Errorf("Key store %s not found or empty", key)
	}

	log.Println(val[1])
	return val[1], nil
}
