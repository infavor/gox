package consulx

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/timer"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

// ConsulClient is a tool for operating with consul server.
type ConsulClient struct {
	// Servers is an array of server:port group
	Servers                   []string
	TTL                       time.Duration
	DeregisterCriticalService bool
	Service                   *api.AgentServiceRegistration
	currentServerIndex        int
	currentApiClient          *api.Client
	check                     *api.AgentServiceCheck
	// check service health
	ServiceCheck func() error
}

// checkServers checks all configured servers's status.
func (client *ConsulClient) checkServers() int {
	log.Debug("checking consul servers...")
	failCount := 0
	for _, s := range client.Servers {
		gox.Try(func() {
			failCount += checkServer(s)
		}, func(i interface{}) {
			failCount++
			log.Error("error when checking server status:", i)
		})
	}
	log.Debug("server check sum, healthy servers ", len(client.Servers)-failCount, ", error servers: ", failCount)
	return failCount
}

// Run runs a consul client.
func (client *ConsulClient) Run() {
	client.currentServerIndex = -1
	if client.Servers == nil || len(client.Servers) == 0 {
		log.Error("no server configured")
		return
	}
	failCount := client.checkServers()
	if failCount == len(client.Servers) {
		log.Error("all consul server is unavailable!")
		return
	}
	client.switchRegisterServer(true)

	client.check = &api.AgentServiceCheck{
		TTL: client.TTL.String(),
	}
	if client.DeregisterCriticalService {
		client.check.DeregisterCriticalServiceAfter = client.TTL.String()
	}

	if client.Service.Check == nil {
		client.Service.Check = client.check
	}

	timer.Start(0, client.TTL/2, 0, func() {
		client.registerService(0)
	})
}

// switchRegisterServer switches consul server if current consul is not available.
func (client *ConsulClient) switchRegisterServer(first bool) {
	if client.currentServerIndex >= len(client.Servers) {
		client.currentServerIndex = -1
	}
	client.currentServerIndex++
	config := api.DefaultConfig()
	config.Address = client.Servers[client.currentServerIndex]
	if !first {
		log.Warn("switch server to ", config.Address)
	}
	c, err := api.NewClient(config)
	if err != nil {
		log.Error("server is not available: ", config.Address)
		return
	}
	client.currentApiClient = c
}

// refresh keeps register status for services.
func (client *ConsulClient) registerService(retry int) {
	log.Debug("refresh services...")
	for client.currentApiClient == nil {
		client.switchRegisterServer(false)
	}
	if err := client.currentApiClient.Agent().ServiceRegister(client.Service); err != nil {
		log.Error("error register service[ID:", client.Service.ID, ", Name:", client.Service.Name, "]: ", err)
		time.Sleep(time.Second * 5)
		retry++
		if retry >= 3 {
			log.Error("register failed finally, retry next round")
			client.currentApiClient = nil
			return
		}
		client.registerService(retry)
	}
	log.Debug("refresh service[ID:", client.Service.ID, ", Name:", client.Service.Name, "] success")
	if err := client.loadServices(); err != nil {
		log.Error("error get services info: ", err)
	}
}

func (client *ConsulClient) loadServices() error {
	ss, err := client.currentApiClient.Agent().Services()
	if err != nil {
		return err
	}
	fmt.Println(ss)
	return nil
}

// checkServer checks if a consul server is available.
func checkServer(server string) int {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}
	_, status, err := httpx.Mock().URL(server).Do()
	if err != nil {
		log.Warn(err)
		return 1
	}
	if err != nil || status != 200 {
		log.Warn("err server status: ", status, " while checking server ", server)
		return 1
	}
	return 0
}
