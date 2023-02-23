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

const SelefraServerHost = "SELEFRA_SERVER_URL"

// GetServerHost 获取服务器的地址
func GetServerHost() string {
	return os.Getenv(SelefraServerHost)
}

const SelefraCloudToken = "SELEFRA_CLOUD_TOKEN"

func GetCloudToken() string {
	return os.Getenv(SelefraCloudToken)
}

// ------------------------------------------------- --------------------------------------------------------------------
