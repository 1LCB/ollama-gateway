package utils

import (
	"errors"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Balancer struct {
	Servers []*Server
	mu      sync.RWMutex
	index   atomic.Uint32
}

type Server struct {
	Address string
	Healthy atomic.Bool
}

func NewLoadBalancer(serverAddresses []string) *Balancer {
	intialServers := make([]*Server, 0, len(serverAddresses))
	for _, address := range serverAddresses {
		intialServers = append(intialServers, &Server{Address: address})
	}

	newLoadBalancer := &Balancer{
		Servers: intialServers,
		index: atomic.Uint32{},
		mu: sync.RWMutex{},
	}

	go newLoadBalancer.HealthChecker()
	return newLoadBalancer
}

func (x *Balancer) SetServers(addresses []string) {
	newServers := make(map[string]bool)
	for _, address := range addresses {
		newServers[address] = false
	}

	x.mu.Lock()
	defer x.mu.Unlock()

	for i := range x.Servers {
		existingServer := x.Servers[i]
		if _, exists := newServers[existingServer.Address]; exists && existingServer.Healthy.Load() {
			newServers[existingServer.Address] = true
		}
	}

	x.Servers = make([]*Server, 0, len(newServers))
	for address, alive := range newServers {
		var server Server
		server.Address = address
		server.Healthy.Store(alive)

		x.Servers = append(x.Servers, &server)
	}
}


func (x *Balancer) HealthChecker() {
	var wg sync.WaitGroup
	client := &http.Client{}
	semaphore := make(chan struct{}, 10)

	for {
		client.Timeout = time.Duration(cfg.LoadBalancer.HealthCheckTimeoutInMillis) * time.Millisecond

		x.mu.RLock()
		for i := range x.Servers {
			server := x.Servers[i]

			semaphore <- struct{}{}

			wg.Add(1)
			go x.checkServer(client, server, &wg, semaphore)
		}
		x.mu.RUnlock()

		wg.Wait()

		time.Sleep(time.Second * time.Duration(cfg.LoadBalancer.HealthCheckIntervalInSeconds))
	}
}


func (*Balancer) checkServer(client *http.Client, server *Server, waitGroup *sync.WaitGroup, semaphore chan struct{}) {
	defer func() {
		waitGroup.Done()
		<-semaphore
	}()

	resp, err := client.Get(server.Address + cfg.LoadBalancer.HealthCheckEndpoint)
	if err == nil {
		resp.Body.Close()
	}

	isHealthy := err == nil && (resp.StatusCode >= 200 && resp.StatusCode < 300)
	server.Healthy.Store(isHealthy)
}

func (x *Balancer) GetServerByRoundRobin() (string, error) {
	x.mu.RLock()
	defer x.mu.RUnlock()

	for i := 0; i < len(x.Servers); i++ {
		currentIndex := x.index.Load()
		server := x.Servers[currentIndex]

		x.index.Store((currentIndex + 1) % uint32(len(x.Servers)))

		if server.Healthy.Load() {
			return server.Address, nil
		}
	}
	return "", errors.New("no healthy backend servers available")
}