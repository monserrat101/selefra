package provider

import (
	"testing"
)

func TestSync(t *testing.T) {
	//global.WorkSpace() = "../../tests/workspace/offline"
	errLogs, _, err := Sync()
	if err != nil {
		t.Error(err)
	}
	if len(errLogs) != 0 {
		t.Error(errLogs)
	}
}
