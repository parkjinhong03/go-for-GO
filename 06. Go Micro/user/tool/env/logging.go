package env

import "os"

func GetForLogging() (env string) {
	env = os.Getenv("ENV")
	switch env {
	case "DEV":
		env = "development"
	case "PROD":
		env = "production"
	default:
		env = "development"
	}
	return
}