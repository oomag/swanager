package command

import (
	"fmt"

	"github.com/dokkur/swanager/core/entities"
	"github.com/dokkur/swanager/core/swarm/node"
)

// NodeList returns node list
type NodeList struct {
	CommonCommand

	OnlyAvailable bool

	responseChan chan<- []entities.Node
}

// NewNodeListCommand create command
func NewNodeListCommand(command NodeList) (NodeList, chan []entities.Node, chan error) {
	response := make(chan []entities.Node, 1)
	err := make(chan error, 1)

	command.errorChan = err
	command.responseChan = response
	return command, response, err
}

// Process deletes command
func (nl NodeList) Process() {

	nodes, err := node.List()
	if err != nil {
		nl.errorChan <- fmt.Errorf("Node list error: %s\n", err.Error())
	}

	if !nl.OnlyAvailable {
		nl.responseChan <- nodes
		return
	}

	result := make([]entities.Node, 0)
	for _, node := range nodes {
		if node.Availability != entities.NodeAvailabilityActive ||
			node.State != entities.NodeStateReady {
			continue
		}

		result = append(result, node)
	}

	nl.responseChan <- result
}
