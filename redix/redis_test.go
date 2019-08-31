package redix_test

import (
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/redix"
	"github.com/hetianyi/gox/uuid"
	"testing"
	"time"
)

func TestGetRedis(t *testing.T) {
	client := redix.InitRedisClient("192.168.245.142:6379", "123123", 2)
	for true {
		if client != nil {
			if err := client.Set(uuid.UUID(), uuid.UUID(), time.Minute).Err(); err != nil {
				logger.Error(err)
			} else {
				logger.Info("add key success")
			}
		}
		time.Sleep(time.Second * 5)
	}
}
