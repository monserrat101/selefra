package init

import (
	"context"
	"github.com/selefra/selefra/cli_ui"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/pkg/cli_runtime"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init [project name]",
		Short: "Prepare your working directory for other commands",
		Long:  "Prepare your working directory for other commands",
		RunE:  initFunc,
	}
	cmd.PersistentFlags().BoolP("force", "f", false, "force overwriting the directory if it is not empty")
	cmd.PersistentFlags().StringP("relevance", "r", "", "associate to selefra cloud project, use only after login")

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func initFunc(cmd *cobra.Command, args []string) error {

	cli_runtime.Init("./")

	downloadDirectory, err := config.GetDefaultDownloadCacheDirectory()
	if err != nil {
		return err
	}

	relevance, _ := cmd.PersistentFlags().GetString("relevance")
	force, _ := cmd.PersistentFlags().GetBool("force")

	dsn, diagnostics := cli_runtime.GetDSN()
	if err := cli_ui.PrintDiagnostics(diagnostics); err != nil {
		return err
	}

	NewInitCommandExecutor(&InitCommandExecutorOptions{
		DownloadWorkspace: downloadDirectory,
		ProjectWorkspace:  "./",
		IsForceInit:       force,
		RelevanceProject:  relevance,
		DSN:               dsn,
	}).Run(context.Background())

	return nil
}

//func initFunc(cmd *cobra.Command, args []string) error {
//	//if err := setInitGlobalVariable(args); err != nil {
//	//	return err
//	//}
//
//	// 1. check if workspace dir exist
//	_, err := os.Stat(global.WorkSpace())
//	if errors.Is(err, os.ErrNotExist) {
//		err = os.Mkdir(global.WorkSpace(), 0755)
//		if err != nil {
//			return nil
//		}
//	}
//	dir, _ := os.ReadDir(global.WorkSpace())
//	for i, v := range dir { // ignore logs dir
//		if v.Name() == "logs" {
//			dir = append(dir[0:i], dir[i+1:]...)
//		}
//	}
//
//	// 2. workspace must be empty or set force flag
//	force, _ := cmd.PersistentFlags().GetBool("force")
//	if len(dir) != 0 && !force {
//		return fmt.Errorf("%s is not empty; Rerun in an empty directory, or use -- force/-f to force overwriting in the current directory\n", global.WorkSpace())
//	}
//
//	// 3. check if workspace is already init
//	if err := reInit(); err != nil {
//		return err
//	}
//
//	// 4. generate yaml file
//	if err := createYaml(cmd); err != nil {
//		//_ = http_client.TrySetUpStage(global.ProjectName(), http_client.Failed)
//	}
//
//	return nil
//}
//
//var storageOptions *postgresql_storage.PostgresqlStorageOptions
//
//func setStorage(ctx context.Context, config *config.SelefraBlock) error {
//	_, diag := pgstorage.PgStorage(ctx)
//	if diag != nil {
//		err := cli_ui.PrintDiagnostic(diag.GetDiagnosticSlice())
//		if err != nil {
//			return errors.New(`The database maybe not ready.
//		You can execute the following command to install the official database image.
//		docker run --name selefra_postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres\n`)
//		}
//	}
//
//	return nil
//}
//
//func setSelefraConfig(cmd *cobra.Command) (*config.SelefraBlock, error) {
//	c := &config.SelefraBlock{}
//	ctx := cmd.Context()
//
//	relevance, _ := cmd.PersistentFlags().GetString("relevance")
//	if relevance != "" {
//		c.Name = relevance
//		if err := login.MustLogin(""); err != nil {
//			return c, errors.New("relevance flag set but can't login")
//		}
//	} else {
//		c.Name = getProjectName()
//		_ = login.ShouldLogin()
//	}
//
//	if err := setStorage(ctx, c); err != nil {
//		return c, err
//	}
//
//	c.CliVersion = version.Version
//	cloudConfig := &config.CloudBlock{
//		Project:      c.Name,
//		Organization: global.OrgName(),
//		HostName:     "",
//	}
//	c.Cloud = cloudConfig
//
//	if err := http_client.TrySetUpStage(c.Name, http_client.Creating); err != nil {
//		return c, err
//	}
//
//	return c, nil
//}
//
//func setProviderConfig(cmd *cobra.Command, configYaml *config.RootConfig) error {
//	ctx := cmd.Context()
//
//	prov, err := getProvidersList()
//	if err != nil {
//		return err
//	}
//
//	provs := cli_ui.SelectProviders(prov)
//	if len(provs) == 0 {
//		return errors.New("No provider selected or user canceled.")
//	}
//	initHeaderOutput(provs)
//
//	namespace, _, err := utils.Home()
//	if err != nil {
//		return err
//	}
//	provider := registry.NewProviderRegistry(namespace)
//
//	for _, s := range provs {
//		pr := registry.Provider{
//			Name:    s,
//			Version: "latest",
//			Source:  "",
//		}
//		p, err := provider.Download(ctx, pr, true)
//		if err != nil {
//			return fmt.Errorf("	Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
//		}
//		cli_ui.Successf("	Installed %s@%s verified\n", p.Name, p.Version)
//		cli_ui.Infof("	Synchronization %s@%s's config...\n", p.Name, p.Version)
//
//		plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
//		if err != nil {
//			return fmt.Errorf("	Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
//		}
//
//		plugProvider := plug.Provider()
//		opt, err := json.Marshal(storageOptions)
//		workspace := global.WorkSpace()
//		initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
//			Workspace: &workspace,
//			Storage: &shard.Storage{
//				Type:           0,
//				StorageOptions: opt,
//			},
//			IsInstallInit:  pointer.TruePointer(),
//			ProviderConfig: pointer.ToStringPointer(""),
//		})
//
//		if err != nil {
//			return err
//		}
//		if initRes != nil && initRes.Diagnostics != nil {
//			err := cli_ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
//			if err != nil {
//				return err
//			}
//		}
//
//		res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
//		if err != nil {
//			return fmt.Errorf("	Synchronization %s@%s's config failed：%s\n", p.Name, p.Version, err.Error())
//		}
//		cli_ui.Successf("	Synchronization %s@%s's config successful\n", p.Name, p.Version)
//		err = tools.AppendProviderDecl(p, configYaml, "latest")
//		if err != nil {
//			return err
//		}
//		err = tools.SetProviderTmpl(res.DefaultConfigTemplate, p, configYaml)
//		if err != nil {
//			return fmt.Errorf("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
//		}
//	}
//
//	return nil
//}
//
//func createYaml(cmd *cobra.Command) error {
//	configYaml := config.RootConfig{}
//
//	selefraConfig, err := setSelefraConfig(cmd)
//	if err != nil {
//		return err
//	}
//	configYaml.Selefra = *selefraConfig
//
//	if err := setProviderConfig(cmd, &configYaml); err != nil {
//		return err
//	}
//
//	return writeToFile(&configYaml)
//}
//
//func writeToFile(configYaml *config.RootConfig) error {
//	// process selefra config depend on whether user login
//	waitStr, err := yaml.Marshal(configYaml)
//	if err != nil {
//		return err
//	}
//	var selefraConfigStr []byte
//	if global.Token() != "" {
//		var initConfigYaml config.RootConfigInitWithLogin
//		err = yaml.Unmarshal(waitStr, &initConfigYaml)
//		if err != nil {
//			return err
//		}
//
//		selefraConfigStr, err = yaml.Marshal(initConfigYaml)
//	} else {
//		var initConfigYaml config.RootConfigInit
//		err = yaml.Unmarshal(waitStr, &initConfigYaml)
//		if err != nil {
//			return err
//		}
//
//		selefraConfigStr, err = yaml.Marshal(initConfigYaml)
//	}
//
//	// create dir for rules
//	rulePath := filepath.Join(global.WorkSpace(), "rules")
//	_, err = os.Stat(rulePath)
//	if err != nil {
//		if os.IsNotExist(err) {
//			mkErr := os.Mkdir(rulePath, 0755)
//			if mkErr != nil {
//				return mkErr
//			}
//		}
//	}
//
//	err = os.WriteFile(filepath.Join(rulePath, "iam_mfa.yaml"), []byte(strings.TrimSpace(ruleComment)), 0644)
//	if err != nil {
//		return err
//	}
//	err = os.WriteFile(filepath.Join(global.WorkSpace(), "module.yaml"), []byte(strings.TrimSpace(moduleComment)), 0644)
//	if err != nil {
//		return err
//	}
//	err = os.WriteFile(filepath.Join(global.WorkSpace(), "selefra.yaml"), selefraConfigStr, 0644)
//
//	cli_ui.Successf(`
//Selefra has been successfully initialized!
//
//Your new Selefra project "%s" was created!
//
//To perform an initial analysis, run selefra apply
//	`, configYaml.Selefra.Name)
//
//	return nil
//}
//
//// reInit check if current workspace is selefra workspace, then tell user to choose if rewrite selefra workspace
//func reInit() error {
//	//_, err := config.GetConfig()
//	//if err != nil && errors.Is(err, config.ErrNotSelefra) {
//	//	return nil
//	//}
//
//	reader := bufio.NewReader(os.Stdin)
//	cli_ui.Warningf("Warning: %s is already init. Continue and overwrite it?[Y/N]", global.WorkSpace())
//	text, err := reader.ReadString('\n')
//	text = strings.TrimSpace(strings.ToLower(text))
//	if err != nil {
//		return nil
//	}
//	if text != "y" && text != "Y" {
//		return errors.New("config file already exists")
//	}
//
//	return nil
//}
//
////func setInitGlobalVariable(args []string) error {
////	// set global workspace
////	wd, err := os.Getwd()
////	if err != nil {
////		return err
////	}
////	dirname := "."
////	if len(args) > 0 {
////		dirname = args[0]
////	}
////	global.Init("init", global.WithWorkspace(filepath.Join(wd, dirname)))
////
////	return nil
////}
//
//func getProvidersList() ([]string, error) {
//	var prov []string
//	cli_ui.Infoln("Getting provider list...")
//	req, _ := http.NewRequest("GET", "https://github.com/selefra/registry/file-list/main/provider", nil)
//	client := http.Client{}
//	res, err := client.Do(req)
//	if err != nil {
//		cli_ui.Errorf("Error: %s", err.Error())
//		return nil, err
//	}
//	d, err := goquery.NewDocumentFromReader(res.Body)
//	if err != nil {
//		cli_ui.Errorf("Error: %s", err.Error())
//		return nil, err
//	}
//	d.Find(".js-navigation-open.Link--primary").Each(func(i int, s *goquery.Selection) {
//		if s.Text() != "template" {
//			prov = append(prov, s.Text())
//		}
//	})
//	return prov, nil
//}
//
//// Generate the name for the current project
//func getProjectName() (projectName string) {
//	defer func() {
//		if projectName == "" {
//			projectName = filepath.Base(global.WorkSpace())
//		}
//		global.SetProjectName(projectName)
//	}()
//	var err error
//	reader := bufio.NewReader(os.Stdin)
//
//	fmt.Printf("project name:(%s)", filepath.Base(global.WorkSpace()))
//
//	projectName, err = reader.ReadString('\n')
//	if err != nil {
//		return ""
//	}
//	projectName = strings.TrimSpace(strings.Replace(projectName, "\n", "", -1))
//
//	return projectName
//}
