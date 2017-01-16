package config

import (
	"flag"
	"os"
)

var (
	// Port of API server
	Port string

	// MongoURL url to connect to mongodb
	MongoURL string

	// DatabaseDriver backend db driver
	DatabaseDriver string

	// DatabaseName is a name of database
	DatabaseName string
)

func init() {
	loadConfigFile()

	flag.StringVar(&Port, "p", "4945", "Api port")
	flag.StringVar(&MongoURL, "m", "mongodb://127.0.0.1:27017/swanager", "Mongodb url")
	flag.StringVar(&DatabaseDriver, "d", "mongo", "Database driver (default: mongo)")
	flag.StringVar(&DatabaseName, "db", "swanager", "Database name (default: swanager)")
	flag.Parse()

	Port = getEnvValue("SWANAGER_PORT", Port)
	DatabaseDriver = getEnvValue("SWANAGER_DB_DRIVER", DatabaseDriver)
	DatabaseName = getEnvValue("SWANAGER_DB_NAME", DatabaseName)
}

func getEnvValue(varName string, currentValue string) string {
	if os.Getenv(varName) != "" {
		return os.Getenv(varName)
	}
	return currentValue
}

func loadConfigFile() {
	// TODO: Read config from file
}
