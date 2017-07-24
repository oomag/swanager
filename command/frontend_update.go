package command

import "github.com/dokkur/swanager/frontend"

// FrontendUpdate command to update frontends
type FrontendUpdate struct {
	CommonCommand
}

// Process service logs
func (fu FrontendUpdate) Process() {
	frontend.Update()
}
