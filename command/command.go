package command

// Command interface
type Command interface {
	Process()
}

// CommonCommand common fields
type CommonCommand struct {
	errorChan chan<- error
}

// RunAsync runs command asynchronously
func RunAsync(command Command) {
	go command.Process()
}
