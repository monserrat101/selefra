package apply

import (
	"context"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"testing"
)

func TestApplyCommandExecutor_Run(t *testing.T) {
	projectWorkspace := "./test_data"
	downloadWorkspace := "./test_download"

	cli_runtime.Init(projectWorkspace)

	NewApplyCommandExecutor(&ApplyCommandExecutorOptions{
		ProjectWorkspace:  projectWorkspace,
		DownloadWorkspace: downloadWorkspace,
	}).Run(context.Background())
}
