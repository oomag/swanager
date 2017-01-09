package config

import (
	"flag"
	"os"
)

var (
	// Port of API server
	Port string
)

func init() {
	flag.StringVar(&Port, "p", "4945", "Api port")

	flag.Parse()

	Port = getEnvValue("SWANAGER_PORT", Port)
}

func getEnvValue(varName string, currentValue string) string {
	if os.Getenv(varName) != "" {
		return os.Getenv(varName)
	}
	return currentValue
}
