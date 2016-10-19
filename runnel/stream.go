package runnel

import (
	"errors"
	"time"
)

// Stream output
func (run *Runnel) Stream(key string, timeout int) (string, error) {

	if key == "" {
		return "", errors.New("Empty key")
	}

	if timeout == 0 {
		timeout = 10
	}

	val, err := run.Redis.BRPop(
		(time.Duration(timeout) * time.Second),
		key).Result()

	if err != nil {
		return "", err
	}

	return val[1], nil
}
