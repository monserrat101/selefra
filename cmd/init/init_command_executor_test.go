package init

import (
	"context"
	"os"
	"testing"
)

func TestInitCommandExecutor_Run(t *testing.T) {

	_ = os.Setenv(SelefraInputInitForceConfirm, "y")

	NewInitCommandExecutor(&InitCommandExecutorOptions{
		DownloadWorkspace: "./test_download",
		ProjectWorkspace:  "./test_data",
		IsForceInit:       true,
		RelevanceProject:  "",
		DSN:               "",
	}).Run(context.Background())
}
