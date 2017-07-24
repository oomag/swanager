package frontend

import (
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/Sirupsen/logrus"
	"github.com/davecgh/go-spew/spew"
	"github.com/dokkur/swanager/config"
	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm"
	"github.com/dokkur/swanager/core/swarm/node"
	vampRouter "github.com/dokkur/swanager/frontend/vamp_router"
)

// Updatable notifies some component with service
type Updatable interface {
	Update([]entities.Service, []entities.Node)
}

var frontends = make([]Updatable, 0)
var lock sync.Mutex

// Init - inits frontend update processing
func Init() {
	spew.Dump(config.VampRouterURL)
	if config.VampRouterURL != "" {
		frontends = append(frontends, &vampRouter.VampRouter{
			URL: config.VampRouterURL,
		})
	}
}

// Update updates frontend config
//   - Only one update can be run at a time, others will be blocked, consider to run it in coroutine
func Update() {
	log().Debugln("Frontend update requested.")
	if len(frontends) == 0 {
		log().Debugln("Nothing to update.")
		return
	}

	// Only one update at a time
	lock.Lock()
	defer lock.Unlock()

	runningServices, err := getRunningServices()
	if err != nil {
		return
	}

	nodes, err := getAvailableNodes()
	if err != nil {
		return
	}

	log().Debugf("Updating %d running services on %d nodes", len(runningServices), len(nodes))

	var wg sync.WaitGroup
	for _, frontend := range frontends {
		wg.Add(1)
		go func(front Updatable, runServs []entities.Service, wg *sync.WaitGroup) {
			defer wg.Done()
			front.Update(runServs, nodes)
		}(frontend, runningServices, &wg)
	}
	wg.Wait()
}

func getAvailableNodes() (nodes []entities.Node, err error) {
	allNodes, err := node.List()
	if err != nil {
		return
	}

	nodes = make([]entities.Node, 0)
	for _, node := range allNodes {
		if node.Availability != entities.NodeAvailabilityActive ||
			node.State != entities.NodeStateReady {
			continue
		}

		nodes = append(nodes, node)
	}
	return
}

func getRunningServices() (services []entities.Service, err error) {
	allServices, err := entities.GetServices(bson.M{})
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for index := range allServices {
		wg.Add(1)
		go func(serv *entities.Service, wg *sync.WaitGroup) {
			defer wg.Done()

			swarm.GetServiceStatuses(serv)
		}(&allServices[index], &wg)
	}
	wg.Wait()

	services = make([]entities.Service, 0)
	for _, service := range allServices {
		if len(service.Status) > 0 && len(service.FrontendEndpoints) > 0 {
			services = append(services, service)
		}
	}
	return
}

func log() *logrus.Entry {
	return logrus.WithField("module", "frontend")
}
