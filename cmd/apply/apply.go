package apply

import (
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/spf13/cobra"
)

func NewApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "apply",
		Short:            "Analyze infrastructure",
		Long:             "Analyze infrastructure",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE:             apply,
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

// ------------------------------------------------- --------------------------------------------------------------------

func apply(cmd *cobra.Command, args []string) error {

	cli_runtime.Init("./")

	projectWorkspace := "./"
	downloadWorkspace, _ := config.GetDefaultDownloadCacheDirectory()
	NewApplyCommandExecutor(&ApplyCommandExecutorOptions{
		ProjectWorkspace:  projectWorkspace,
		DownloadWorkspace: downloadWorkspace,
	}).Run(cmd.Context())
	return nil
}

// ------------------------------------------------- --------------------------------------------------------------------
