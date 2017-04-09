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
	//DatabaseDriver string

	// DatabaseName is a name of database
	DatabaseName string

	// MountPathPrefix base mounted share path
	MountPathPrefix string

	// LogFileName path to logfile
	LogFileName string

	// LocalSecretKey secret key to authenticate local services
	LocalSecretKey string
)

func init() {
	loadConfigFile()

	flag.StringVar(&Port, "p", "4945", "Api port")
	flag.StringVar(&LogFileName, "l", "", "Path to log file (default: stdout)")
	flag.StringVar(&MongoURL, "m", "mongodb://127.0.0.1:27017/swanager", "Mongodb url")
	//flag.StringVar(&DatabaseDriver, "d", "mongo", "Database driver (default: mongo)")
	flag.StringVar(&DatabaseName, "db", "swanager", "Database name (default: swanager)")
	flag.StringVar(&MountPathPrefix, "share", "/data", "Mount point base path (default: /data)")
	flag.StringVar(&LocalSecretKey, "lsk", "", "Secret key to authenticate local services (default: none, won't be authenticated)")
	flag.Parse()

	Port = getEnvValue("SWANAGER_PORT", Port)
	LogFileName = getEnvValue("SWANAGER_LOG", LogFileName)
	MongoURL = getEnvValue("SWANAGER_MONGO_URL", MongoURL)
	//DatabaseDriver = getEnvValue("SWANAGER_DB_DRIVER", DatabaseDriver)
	DatabaseName = getEnvValue("SWANAGER_DB_NAME", DatabaseName)
	MountPathPrefix = getEnvValue("SWANAGER_PATH_PREFIX", MountPathPrefix)
	LocalSecretKey = getEnvValue("SWANAGER_LOCAL_SECRET_KEY", LocalSecretKey)
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
