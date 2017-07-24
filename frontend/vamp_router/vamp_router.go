package vampRouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/magneticio/vamp-router/haproxy"

	"github.com/dokkur/swanager/core/entities"
)

// VampRouter represenst vamp-router integration
type VampRouter struct {
	URL string

	nodes []entities.Node

	frontendsMutex sync.Mutex
	frontends      map[uint32]haproxy.Frontend

	backendsMutex sync.Mutex
	backends      []*haproxy.Backend
}

// Update updates vamp-router configuration
func (vr *VampRouter) Update(services []entities.Service, nodes []entities.Node) {
	vr.cleanup()
	vr.nodes = nodes

	if len(services) == 0 {
		return
	}

	var wg sync.WaitGroup
	for _, service := range services {
		if len(service.FrontendEndpoints) == 0 {
			continue
		}

		wg.Add(1)
		go func(serv entities.Service, wg *sync.WaitGroup) {
			defer wg.Done()
			vr.parseService(serv)
		}(service, &wg)
	}
	wg.Wait()
	// spew.Dump(vr.frontends)
	// spew.Dump(vr.backends)

	frontends := make([]*haproxy.Frontend, 0, len(vr.frontends))
	for _, front := range vr.frontends {
		frontends = append(frontends, &front)
	}

	config := haproxy.Config{
		Backends:  vr.backends,
		Frontends: frontends,
		Routes:    make([]haproxy.Route, 0),
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		return
	}

	// fmt.Println(string(configJSON))

	http.Post(vr.URL+"/v1/config", "application/json", bytes.NewReader(configJSON))
}

func (vr *VampRouter) parseService(service entities.Service) {
	for _, endpoint := range service.FrontendEndpoints {
		vr.parseEndpoint(service, endpoint)
	}
}

func (vr *VampRouter) parseEndpoint(service entities.Service,
	endpoint entities.FrontendEndpoint) {

	backendName := fmt.Sprintf("swarm_%s_%d", service.NSName, endpoint.InternalPort)
	filterName := fmt.Sprintf("front_%s_%d", service.NSName, endpoint.ExternalPort)

	vr.addFrontendFilter(endpoint.ExternalPort, haproxy.Filter{
		Name:        filterName,
		Destination: backendName,
		Condition:   fmt.Sprintf("hdr(Host) -i %s", endpoint.Domain),
	})

	vr.addBackend(haproxy.Backend{
		Name:      backendName,
		Mode:      "http",
		ProxyMode: true,
		Options: haproxy.ProxyOptions{
			ForwardFor: true,
			HttpCheck:  true,
		},
		Servers: vr.getServers(endpoint.InternalPort),
	})
}

func (vr *VampRouter) getServers(port uint32) []*haproxy.ServerDetail {
	servers := make([]*haproxy.ServerDetail, 0)

	for _, node := range vr.nodes {
		servers = append(servers, &haproxy.ServerDetail{
			Name:          node.ID,
			Host:          node.Addr,
			Port:          int(port),
			Check:         true,
			CheckInterval: 10,
		})
	}

	return servers
}

func (vr *VampRouter) addFrontendFilter(port uint32, filter haproxy.Filter) {
	vr.frontendsMutex.Lock()
	defer vr.frontendsMutex.Unlock()

	if _, ok := vr.frontends[port]; !ok {

		vr.frontends[port] = haproxy.Frontend{
			Name:           fmt.Sprintf("port%d", port),
			BindIp:         "0.0.0.0",
			BindPort:       int(port),
			DefaultBackend: "default",
			Mode:           "http",
			Options: haproxy.ProxyOptions{
				HttpClose: true,
			},
		}
	}

	front := vr.frontends[port]
	front.Filters = append(front.Filters, &filter)

	vr.frontends[port] = front
}

func (vr *VampRouter) addBackend(back haproxy.Backend) {
	vr.backendsMutex.Lock()
	defer vr.backendsMutex.Unlock()

	vr.backends = append(vr.backends, &back)
}

func (vr *VampRouter) cleanup() {
	vr.frontendsMutex.Lock()
	defer vr.frontendsMutex.Unlock()
	vr.backendsMutex.Lock()
	defer vr.backendsMutex.Unlock()

	vr.frontends = make(map[uint32]haproxy.Frontend, 0)
	vr.backends = make([]*haproxy.Backend, 0)
}
