package logs

import (
	"bufio"
	"context"
	"fmt"

	"github.com/dokkur/swanager/core/entities"
	swarm_service "github.com/dokkur/swanager/core/swarm/service"
)

var ctx context.Context
var cancel context.CancelFunc

// Start - Starts logs listening service
func Start() {
	ctx, cancel = context.WithCancel(context.Background())
}

// For listens for service logs
func For(service *entities.Service) {
	ch := make(chan string)

	go func(ch chan string) {
		reader, err := swarm_service.LogsFollow(service, 20)
		if err != nil {
			return
		}

		scanner := bufio.NewScanner(reader)
		result := make([]string, 0)

		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		for {
			ch <- scanner.Text()
		}
	}(ch)

	listening := true
	for listening {
		select {
		case stdin, ok := <-ch:
			if !ok {
				listening = false
				break
			} else {
				fmt.Println("Read input from stdin:", stdin)
			}
			// case <-time.After(1 * time.Second):
			// Do something when there is nothing read from stdin
		}
	}
}

// Stop - Stop service
func Stop() {
	cancel()
}
