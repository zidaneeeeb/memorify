package environment

import "os"

// ServiceEnv denotes environment in which a service is
// running.
type ServiceEnv string

// Followings are the known ServiceEnv values.
const (
	DevelopmentEnv ServiceEnv = "development"
	StagingEnv     ServiceEnv = "staging"
	ProductionEnv  ServiceEnv = "production"
)

// serviceEnvKey is the environment variable key in which
// service environment value is stored.
const serviceEnvKey string = "BASEENV"

// ServiceEnv return TKPENV service environment.
func GetServiceEnv() ServiceEnv {
	e := os.Getenv(serviceEnvKey)
	if e == "" {
		return DevelopmentEnv
	}
	return ServiceEnv(e)
}
