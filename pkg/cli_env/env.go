package cli_env

import (
	"os"
	"strings"
)

const SelefraCloudFlag = "SELEFRA_CLOUD_FLAG"

// IsCloudEnv Check whether the system is running in the cloud environment
func IsCloudEnv() bool {
	flag := strings.ToLower(os.Getenv(SelefraCloudFlag))
	return flag == "true" || flag == "enable"
}

//func GetTaskID() {
//
//}

// ------------------------------------------------- --------------------------------------------------------------------

const SelefraServerHost = "SELEFRA_CLOUD_HOST"

// GetServerHost Gets the address of the server
func GetServerHost() string {
	return os.Getenv(SelefraServerHost)
}

const SelefraCloudToken = "SELEFRA_CLOUD_TOKEN"

func GetCloudToken() string {
	return os.Getenv(SelefraCloudToken)
}

// ------------------------------------------------- --------------------------------------------------------------------
