package ws

import "github.com/da4nik/swanager/core/entities"

// Notify notify user about state change
func Notify(service *entities.Service) {
	for _, client := range clients {
		client <- *service
	}
}
