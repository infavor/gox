package consulx_test

import (
	"github.com/hashicorp/consul/api"
	"github.com/hetianyi/gox/consulx"
	"testing"
	"time"
)

func TestConsulClient(t *testing.T) {
	var client = &consulx.ConsulClient{
		Servers:                   []string{"192.168.25.132:8500"},
		TTL:                       time.Minute,
		DeregisterCriticalService: true,
		Service: &api.AgentServiceRegistration{
			ID:                "storage-1", // 服务的唯一ID(单个实例)
			Kind:              api.ServiceKindTypical,
			Name:              "godfs-storage", // 服务的名称
			Port:              8076,
			Address:           "192.168.0.123",
			Tags:              []string{"storage", "west"}, // 服务的标签
			EnableTagOverride: true,
		},
	}
	client.Run()

	wait := make(chan int)
	<-wait
}
