package init

import (
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra/config"
	"github.com/spf13/cobra"
	"os"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project name]",
		Short: "Prepare your working directory for other commands",
		Long:  "Prepare your working directory for other commands",
		RunE: func(cmd *cobra.Command, args []string) error {

			relevance, _ := cmd.PersistentFlags().GetString("relevance")
			force, _ := cmd.PersistentFlags().GetBool("force")

			downloadDirectory, err := config.GetDefaultDownloadCacheDirectory()
			if err != nil {
				return err
			}

			projectWorkspace := "./"

			dsn := os.Getenv(env.DatabaseDsn)

			return NewInitCommandExecutor(&InitCommandExecutorOptions{
				IsForceInit:       force,
				RelevanceProject:  relevance,
				ProjectWorkspace:  projectWorkspace,
				DownloadWorkspace: downloadDirectory,
				DSN:               dsn,
			}).Run(cmd.Context())
		},
	}
	cmd.PersistentFlags().BoolP("force", "f", false, "force overwriting the directory if it is not empty")
	cmd.PersistentFlags().StringP("relevance", "r", "", "associate to selefra cloud project, use only after login")

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}
