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

const DefaultCloudHost = "main-grpc.selefra.io"

// GetServerHost Gets the address of the server
func GetServerHost() string {

	// read from env
	if os.Getenv(SelefraServerHost) != "" {
		return os.Getenv(SelefraServerHost)
	}

	return DefaultCloudHost
}

// ------------------------------------------------- --------------------------------------------------------------------

const SelefraCloudToken = "SELEFRA_CLOUD_TOKEN"

func GetCloudToken() string {
	return os.Getenv(SelefraCloudToken)
}

// ------------------------------------------------- --------------------------------------------------------------------
