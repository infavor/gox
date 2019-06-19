package consulx_test

import (
	"github.com/hashicorp/consul/api"
	"github.com/hetianyi/gox/consulx"
	"github.com/hetianyi/gox/logger"
	"testing"
	"time"
)

func init() {
	logger.Init(&logger.Config{
		Level: logger.TraceLevel,
	})
}

func TestConsulClient(t *testing.T) {
	var client = &consulx.ConsulClient{
		Servers:                   []string{"192.168.0.104:8500", "192.168.0.105:8500"},
		TTL:                       time.Minute,
		DeregisterCriticalService: true,
		Service: &api.AgentServiceRegistration{
			ID:                "storage-1", // 服务的唯一ID(单个实例)
			Kind:              api.ServiceKindTypical,
			Name:              "godfs-storage", // 服务的名称
			Port:              8076,
			Address:           "192.168.0.101",
			Tags:              []string{"storage", "west"}, // 服务的标签
			EnableTagOverride: true,
		},
	}
	client.Run()
	logger.Info("consul set passed..................")
	wait := make(chan int)
	<-wait
}
