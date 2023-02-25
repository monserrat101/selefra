package query

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-provider-sdk/storage_factory"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
)

func NewQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "query",
		Short:            "Query infrastructure data from pgstorage",
		Long:             "Query infrastructure data from pgstorage",
		PersistentPreRun: global.DefaultWrappedInit(),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			cli_runtime.Init("./")

			cli_ui.Warningln("Please select table.")

			dsn, d := cli_runtime.GetDSN()
			if utils.HasError(d) {
				_ = cli_ui.PrintDiagnostics(d)
				return
			}
			options := postgresql_storage.NewPostgresqlStorageOptions(dsn)
			storage, diagnostics := storage_factory.NewStorage(context.Background(), storage_factory.StorageTypePostgresql, options)
			if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
				return
			}

			queryClient, _ := NewQueryClient(ctx, storage_factory.StorageTypePostgresql, storage)
			queryClient.Run(ctx)

		},
	}
	return cmd
}
