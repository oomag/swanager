package command

// Command interface
type Command interface {
	Process()
}

// RunAsync runs command asynchronously
func RunAsync(command Command) {
	go command.Process()
}
