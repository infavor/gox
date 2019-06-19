package consulx

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"strings"
	"time"
)

// ConsulClient is a tool for operating with consul server.
type ConsulClient struct {
	Servers                   []string // Servers is an array of server:port group
	TTL                       time.Duration
	DeregisterCriticalService bool
	Service                   *api.AgentServiceRegistration
	currentServerIndex        int
	currentApiClient          *api.Client
	check                     *api.AgentServiceCheck
	ServiceCheck              func() error // check service health
	renewLock                 chan byte
}

// checkServers checks all configured servers's status.
func (client *ConsulClient) checkServers() int {
	logger.Debug("checking consul servers...")
	failCount := 0
	for _, s := range client.Servers {
		gox.Try(func() {
			failCount += checkServer(s)
		}, func(i interface{}) {
			failCount++
			logger.Error("error when checking server status:", i)
		})
	}
	logger.Debug("server check sum, healthy servers ", len(client.Servers)-failCount, ", error servers: ", failCount)
	return failCount
}

// Run runs a consul client.
func (client *ConsulClient) Run() {
	client.currentServerIndex = -1
	if client.Servers == nil || len(client.Servers) == 0 {
		logger.Error("no server configured")
		return
	}
	failCount := client.checkServers()
	if failCount == len(client.Servers) {
		logger.Error("all consul server is unavailable!")
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

	client.renewLock = make(chan byte)
	go func() {
		for {
			client.renew()
		}
	}()
	client.registerService()
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
		logger.Warn("switch server to ", config.Address)
	}
	c, err := api.NewClient(config)
	if err != nil {
		logger.Error("server is not available: ", config.Address)
		return
	}
	client.currentApiClient = c
}

// refresh keeps register status for services.
func (client *ConsulClient) registerService() {
	retry := 0
	for {
		logger.Debug("register service...")
		for client.currentApiClient == nil {
			client.switchRegisterServer(false)
		}
		err := client.currentApiClient.Agent().ServiceRegister(client.Service)
		if err != nil {
			logger.Error("error register service[ID:", client.Service.ID, ", Name:", client.Service.Name, "]: ", err)
			time.Sleep(time.Second * 15)
			retry++
			if retry%3 == 0 {
				client.currentApiClient = nil
			}
			continue
		}
		logger.Info("register service success: ", client.Service.Name)
		break
	}
	client.renewLock <- 0
	if err := client.loadServices(); err != nil {
		logger.Error("error get services info: ", err)
	}
}

func (client *ConsulClient) renew() {
	<-client.renewLock
	timer.Start(0, client.TTL/2, 0, func(t *timer.Timer) {
		logger.Debug("try to renew service: ", client.Service.Name)
		err := client.currentApiClient.Agent().PassTTL("service:"+gox.TValue(client.Service.ID == "", client.Service.Name, client.Service.ID).(string), "")
		if err != nil {
			logger.Error("error renew service: ", client.Service.Name)
			t.Destroy()
			client.registerService()
			return
		}
		logger.Debug("renew service[ID:", client.Service.ID, ", Name:", client.Service.Name, "] success")
	})
}

func (client *ConsulClient) loadServices() error {
	ss, err := client.currentApiClient.Agent().Services()
	if err != nil {
		return err
	}
	bs, _ := json.MarshalIndent(ss, "", " ")
	fmt.Println(string(bs))
	return nil
}

// checkServer checks if a consul server is available.
func checkServer(server string) int {
	if !strings.HasPrefix(server, "http://") && !strings.HasPrefix(server, "https://") {
		server = "http://" + server
	}
	_, status, err := httpx.Mock().URL(server).Do()
	if err != nil {
		logger.Warn(err)
		return 1
	}
	if err != nil || status != 200 {
		logger.Warn("err server status: ", status, " while checking server ", server)
		return 1
	}
	return 0
}
