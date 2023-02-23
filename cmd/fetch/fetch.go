package fetch

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/executors"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/module_loader"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/spf13/cobra"
)

func NewFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "fetch",
		Short:            "Fetch resources from configured providers",
		Long:             "Fetch resources from configured providers",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, args []string) error {
			//errFlag := false
			//ctx := cmd.Context()
			//rootConfig, err := config.GetConfig()
			//if err != nil {
			//	return err
			//}
			//cli_ui.Successf("Selefra start fetch")
			//for _, p := range rootConfig.Selefra.ProviderDecls {
			//	prvds := tools.ProvidersByID(rootConfig, p.Name)
			//	for _, prvd := range prvds {
			//		err = Fetch(ctx, p, prvd)
			//		if err != nil {
			//			cli_ui.Errorln(err.Error())
			//			errFlag = true
			//			return err
			//		}
			//	}
			//}
			//if errFlag {
			//	cli_ui.Errorf(`
			//This may be exception, view detailed exception in %s.`,
			//		filepath.Join(global.WorkSpace(), "logs"))
			//
			//}
			//return nil
			//
			//// TODO It's pretty much the same as the first half of apply
			//
			//return nil

			//projectWorkspace := "./test_data/test_fetch_module"
			//downloadWorkspace := "./test_download"

			projectWorkspace := "./"
			downloadWorkspace, _ := config.GetDefaultDownloadCacheDirectory()

			cli_runtime.Init(projectWorkspace)

			Fetch(projectWorkspace, downloadWorkspace)

			return nil
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func Fetch(projectWorkspace, downloadWorkspace string) *schema.Diagnostics {

	// Load the module used for the test
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	loader, err := module_loader.NewLocalDirectoryModuleLoader(&module_loader.LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: &module_loader.ModuleLoaderOptions{
			MessageChannel: messageChannel,
		},
		ModuleDirectory: projectWorkspace,
	})
	if err != nil {
		cli_ui.Errorln(fmt.Sprintf("create local directory module loader from %s error: %s", projectWorkspace, err.Error()))
		return nil
	}
	rootModule, b := loader.Load(context.Background())
	messageChannel.ReceiverWait()
	if !b {
		cli_ui.Errorln(fmt.Sprintf("local directory module loader load  %s failed.", projectWorkspace))
		return nil
	}

	// check module
	validatorContext := module.NewValidatorContext()
	d := rootModule.Check(rootModule, validatorContext)
	if err := cli_ui.PrintDiagnostics(d); err != nil {
		return nil
	}

	// Make an installation plan
	providersInstallPlan, diagnostics := planner.MakeProviderInstallPlan(context.Background(), rootModule)
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return nil
	}
	if len(providersInstallPlan) == 0 {
		cli_ui.Errorln("not providers")
		return nil
	}

	// Installation-dependent dependency
	messageChannel = message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	executor, diagnostics := executors.NewProviderInstallExecutor(&executors.ProviderInstallExecutorOptions{
		Plans:             providersInstallPlan,
		MessageChannel:    messageChannel,
		DownloadWorkspace: downloadWorkspace,
	})
	_ = cli_ui.PrintDiagnostics(diagnostics)
	if utils.HasError(diagnostics) {
		return nil
	}
	d = executor.Execute(context.Background())
	messageChannel.ReceiverWait()
	_ = cli_ui.PrintDiagnostics(d)
	if utils.HasError(d) {
		return nil
	}

	// Develop a data pull plan
	providerFetchPlans, d := planner.NewProviderFetchPlanner(rootModule, providersInstallPlan.ToMap()).MakePlan(context.Background())
	_ = cli_ui.PrintDiagnostics(d)
	if utils.HasError(d) {
		return nil
	}

	// Ready to start pulling
	localProviderManager, err := local_providers_manager.NewLocalProvidersManager(downloadWorkspace)
	if err != nil {
		cli_ui.Errorln("")
		return nil
	}
	messageChannel = message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			_ = cli_ui.PrintDiagnostics(message)
		}
	})
	fetchExecutor := executors.NewProviderFetchExecutor(&executors.ProviderFetchExecutorOptions{
		LocalProviderManager: localProviderManager,
		Plans:                providerFetchPlans,
		MessageChannel:       messageChannel,
		WorkerNum:            3,
		Workspace:            projectWorkspace,
		// TODO
		DSN: env.GetDatabaseDsn(),
	})
	d = fetchExecutor.Execute(context.Background())
	messageChannel.ReceiverWait()
	_ = cli_ui.PrintDiagnostics(d)
	if utils.HasError(d) {
		return nil
	}

	cli_ui.Infoln("fetch done!")

	return nil
}

//func Fetch(ctx context.Context, decl *config.RequireProvider, prvd *config.ProviderBlock) error {
//	if decl.Path == "" {
//		decl.Path = utils.GetPathBySource(*decl.Source, decl.Version)
//	}
//	var providersName = *decl.Source
//	cli_ui.Successf("%s %s@%s pull infrastructure data:\n", prvd.Name, providersName, decl.Version)
//	cli_ui.Print(fmt.Sprintf("Pulling %s@%s Please wait for resource information ...", providersName, decl.Version), false)
//	plug, err := plugin.NewManagedPlugin(decl.Path, providersName, decl.Version, "", nil)
//	if err != nil {
//		return err
//	}
//
//	storageOpt := pgstorage.DefaultPgStorageOpts()
//	pgstorage.WithSearchPath(config.GetSchemaKey(decl, *prvd))(storageOpt)
//
//	opt, err := json.Marshal(storageOpt)
//	if err != nil {
//		return err
//	}
//
//	prvdByte, err := yaml.Marshal(prvd)
//	if err != nil {
//		return err
//	}
//
//	plugProvider := plug.Provider()
//	initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
//		Workspace: pointer.ToStringPointer(global.WorkSpace()),
//		Storage: &shard.Storage{
//			Type:           0,
//			StorageOptions: opt,
//		},
//		IsInstallInit:  pointer.FalsePointer(),
//		ProviderConfig: pointer.ToStringPointer(string(prvdByte)),
//	})
//	if err != nil {
//		return err
//	} else {
//		if initRes.Diagnostics != nil {
//			err := cli_ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
//			if err != nil {
//				return errors.New("fetch plugProvider init error")
//			}
//		}
//	}
//
//	defer plug.Close()
//	dropRes, err := plugProvider.DropTableAll(ctx, &shard.ProviderDropTableAllRequest{})
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return err
//	}
//	if dropRes.Diagnostics != nil {
//		err := cli_ui.PrintDiagnostic(dropRes.Diagnostics.GetDiagnosticSlice())
//		if err != nil {
//			return errors.New("fetch plugProvider drop table error")
//		}
//	}
//
//	createRes, err := plugProvider.CreateAllTables(ctx, &shard.ProviderCreateAllTablesRequest{})
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return err
//	}
//	if createRes.Diagnostics != nil {
//		err := cli_ui.PrintDiagnostic(createRes.Diagnostics.GetDiagnosticSlice())
//		if err != nil {
//			return errors.New("fetch plugProvider create table error")
//		}
//	}
//	var tables []string
//	resources := prvd.Resources
//	if len(resources) == 0 {
//		tables = append(tables, "*")
//	} else {
//		for i := range resources {
//			tables = append(tables, resources[i])
//		}
//	}
//	var maxGoroutines uint64 = 100
//	if prvd.MaxGoroutines > 0 {
//		maxGoroutines = prvd.MaxGoroutines
//	}
//	recv, err := plugProvider.PullTables(ctx, &shard.PullTablesRequest{
//		Tables:        tables,
//		MaxGoroutines: maxGoroutines,
//		Timeout:       0,
//	})
//	if err != nil {
//		cli_ui.Errorln(err.Error())
//		return err
//	}
//	progbar := progress.DefaultProgress()
//	progbar.Add(decl.Name+"@"+decl.Version, -1)
//	success := 0
//	errorsN := 0
//	var total int64
//	for {
//		res, err := recv.Recv()
//		if err != nil {
//			if errors.Is(err, io.EOF) {
//				progbar.Current(decl.Name+"@"+decl.Version, total, "Done")
//				progbar.Done(decl.Name + "@" + decl.Version)
//				break
//			}
//			return err
//		}
//		progbar.SetTotal(decl.Name+"@"+decl.Version, int64(res.TableCount))
//		progbar.Current(decl.Name+"@"+decl.Version, int64(len(res.FinishedTables)), res.Table)
//		total = int64(res.TableCount)
//		if res.Diagnostics != nil {
//			if res.Diagnostics.HasError() {
//				cli_ui.SaveLogToDiagnostic(res.Diagnostics.GetDiagnosticSlice())
//			}
//		}
//		success = len(res.FinishedTables)
//		errorsN = 0
//	}
//	progbar.ReceiverWait(decl.Name + "@" + decl.Version)
//	if errorsN > 0 {
//		cli_ui.Errorf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
//		return nil
//	}
//	cli_ui.Successf("\nPull complete! Total Resources pulled:%d        Errors: %d\n", success, errorsN)
//	return nil
//}
