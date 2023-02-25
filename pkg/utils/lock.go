package utils

import (
	"github.com/selefra/selefra-utils/pkg/id_util"
	"os"
)

// BuildOwnerId The current host name is placed in the owner of the lock so that it is easy to identify who is holding the lock
func BuildOwnerId() string {
	hostname, err := os.Hostname()
	id := id_util.RandomId()
	if err != nil {
		return id
	} else {
		return hostname + "-" + id
	}
}
