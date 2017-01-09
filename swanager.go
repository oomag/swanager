package main

import (
	"os"
	"os/signal"

	_ "github.com/da4nik/swanager/api"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
