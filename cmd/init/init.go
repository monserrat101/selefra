package init

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra/cmd/login"
	"github.com/selefra/selefra/cmd/tools"
	"github.com/selefra/selefra/pkg/httpClient"
	"gopkg.in/yaml.v3"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/selefra/selefra-provider-sdk/grpc/shard"
	"github.com/selefra/selefra-provider-sdk/storage/database_storage/postgresql_storage"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/cmd/version"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/plugin"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/selefra/selefra/ui/term"
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

var storage *postgresql_storage.PostgresqlStorageOptions

func setSelefraConfig(cmd *cobra.Command) (*config.Config, error) {
	c := &config.Config{}
	ctx := cmd.Context()

	// set version
	c.CliVersion = version.Version

	// set db
	storage = postgresql_storage.NewPostgresqlStorageOptions(c.GetDSN())

	_, diag := postgresql_storage.NewPostgresqlStorage(ctx, storage)
	if diag != nil {
		err := ui.PrintDiagnostic(diag.GetDiagnosticSlice())
		if err != nil {
			return c, errors.New(`The database maybe not ready.
		You can execute the following command to install the official database image.
		docker run --name selefra_postgres -p 5432:5432 -e POSTGRES_PASSWORD=pass -d postgres\n`)
		}
	}

	// judge if sync to cloud is need
	var err error
	relevance, _ := cmd.PersistentFlags().GetString("relevance")
	if relevance == "" { // no relevance, try login and setup stage
		c.Name = getProjectName()
		err = login.ShouldLogin("")
		if err != nil { // login failed, return with no error
			return c, nil
		}
	} else { // relevance was set, must login
		c.Name = relevance
		err = login.MustLogin("")
		if err != nil { // login failed, return with an error
			return nil, err
		}
	}

	cloudConfig := &config.Cloud{
		Project:      c.Name,
		Organization: global.ORGNAME,
		HostName:     "",
	}
	c.Cloud = cloudConfig

	if err = httpClient.SetUpStage(global.LOGINTOKEN, c.Name, httpClient.Creating); err != nil {
		return c, err
	}

	return c, nil
}

func setProviderConfig(cmd *cobra.Command, configYaml *config.SelefraConfig) error {
	ctx := cmd.Context()

	prov, err := getProvidersList()
	if err != nil {
		return err
	}

	provs := term.SelectProviders(prov)
	if len(provs) == 0 {
		return errors.New("No provider selected or user canceled.")
	}
	initHeaderOutput(provs)

	token, err := utils.GetCredentialsToken()

	namespace, _, err := utils.Home()
	if err != nil {
		return err
	}
	provider := registry.NewProviderRegistry(namespace)

	for _, s := range provs {
		pr := registry.Provider{
			Name:    s,
			Version: "latest",
			Source:  "",
		}
		p, err := provider.Download(ctx, pr, true)
		if err != nil {
			return fmt.Errorf("	Installed %s@%s failed：%s", p.Name, p.Version, err.Error())
		} else {
			if token != "" && err == nil {
				_ = httpClient.SetUpStage(token, configYaml.Selefra.Cloud.Project, httpClient.Failed)
			}
			ui.PrintSuccessF("	Installed %s@%s verified", p.Name, p.Version)
		}
		ui.PrintInfoF("	Synchronization %s@%s's config...", p.Name, p.Version)
		plug, err := plugin.NewManagedPlugin(p.Filepath, p.Name, p.Version, "", nil)
		if err != nil {
			return fmt.Errorf("	Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}

		plugProvider := plug.Provider()
		opt, err := json.Marshal(storage)
		initRes, err := plugProvider.Init(ctx, &shard.ProviderInitRequest{
			Workspace: global.WORKSPACE,
			Storage: &shard.Storage{
				Type:           0,
				StorageOptions: opt,
			},
			IsInstallInit:  pointer.TruePointer(),
			ProviderConfig: pointer.ToStringPointer(""),
		})

		if err != nil {
			return err
		}
		if initRes != nil && initRes.Diagnostics != nil {
			err := ui.PrintDiagnostic(initRes.Diagnostics.GetDiagnosticSlice())
			if err != nil {
				return err
			}
		}

		res, err := plugProvider.GetProviderInformation(ctx, &shard.GetProviderInformationRequest{})
		if err != nil {
			return fmt.Errorf("	Synchronization %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}
		ui.PrintSuccessF("	Synchronization %s@%s's config successful", p.Name, p.Version)
		err = tools.SetSelefraProvider(p, configYaml, "latest")
		if err != nil {
			return err
		}
		err = tools.SetProviders(res.DefaultConfigTemplate, p, configYaml)
		if err != nil {
			return fmt.Errorf("set %s@%s's config failed：%s", p.Name, p.Version, err.Error())
		}
	}

	return nil
}

func createYaml(cmd *cobra.Command) (*config.SelefraConfig, error) {
	configYaml := config.SelefraConfig{}

	selefraConfig, err := setSelefraConfig(cmd)
	if err != nil {
		return &configYaml, err
	}
	configYaml.Selefra = *selefraConfig

	if err := setProviderConfig(cmd, &configYaml); err != nil {
		return &configYaml, err
	}

	waitStr, err := yaml.Marshal(configYaml)
	if err != nil {
		return &configYaml, err
	}
	var str []byte
	if global.LOGINTOKEN != "" {
		var initConfigYaml config.SelefraConfigInitWithLogin
		err = yaml.Unmarshal(waitStr, &initConfigYaml)
		if err != nil {
			return &configYaml, err
		}

		str, err = yaml.Marshal(initConfigYaml)
	} else {
		var initConfigYaml config.SelefraConfigInit
		err = yaml.Unmarshal(waitStr, &initConfigYaml)
		if err != nil {
			return &configYaml, err
		}

		str, err = yaml.Marshal(initConfigYaml)
	}

	rulePath := filepath.Join(*global.WORKSPACE, "rules")
	_, err = os.Stat(rulePath)
	if err != nil {
		if os.IsNotExist(err) {
			mkErr := os.Mkdir(rulePath, 0755)
			if mkErr != nil {
				return &configYaml, mkErr
			}
		}
	}

	err = os.WriteFile(filepath.Join(rulePath, "iam_mfa.yaml"), []byte(strings.TrimSpace(ruleComment)), 0644)
	if err != nil {
		return &configYaml, err
	}

	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "module.yaml"), []byte(strings.TrimSpace(moduleComment)), 0644)
	if err != nil {
		return &configYaml, err
	}
	err = os.WriteFile(filepath.Join(*global.WORKSPACE, "selefra.yaml"), str, 0644)

	ui.PrintSuccessF(`
Selefra has been successfully initialized! 
	
Your new Selefra project "%s" was created!

To perform an initial analysis, run selefra apply
	`, configYaml.Selefra.Name)

	return &configYaml, nil
}

func initFunc(cmd *cobra.Command, args []string) error {
	if err := setInitGlobalVariable(args); err != nil {
		return err

	}
	// check if workspace dir exist
	_, err := os.Stat(*global.WORKSPACE)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(*global.WORKSPACE, 0755)
		if err != nil {
			return nil
		}
	}
	dir, _ := os.ReadDir(*global.WORKSPACE)

	// ignore logs dir
	for i, v := range dir {
		if v.Name() == "logs" {
			dir = append(dir[0:i], dir[i+1:]...)
		}
	}

	// workspace must be empty or set force flag
	force, _ := cmd.PersistentFlags().GetBool("force")
	if len(dir) != 0 && !force {
		return fmt.Errorf("%s is not empty; Rerun in an empty directory, or use -- force/-f to force overwriting in the current directory\n", *global.WORKSPACE)
	}

	// when relevance was set, user must be login
	relevance, _ := cmd.PersistentFlags().GetString("relevance")
	if relevance != "" {
		if err := login.MustLogin(""); err != nil {
			return errors.New("relevance flag set but can't login")
		}
	}

	// check if workspace is already init
	_, clientErr := config.GetClientStr()
	if !errors.Is(clientErr, config.NoClient) { // workspace already init by selefra
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Error:%s is already init. Continue and overwrite it?[Y/N]\n", *global.WORKSPACE)
		text, err := reader.ReadString('\n')
		text = strings.TrimSpace(strings.ToLower(text))
		if err != nil {
			return nil
		}
		if text != "y" && text != "Y" {
			return errors.New("config file already exists")
		}
	}

	cof, err := createYaml(cmd)
	if err == nil && global.LOGINTOKEN != "" {
		_ = httpClient.SetUpStage(global.LOGINTOKEN, cof.Selefra.Cloud.Project, httpClient.Failed)
	}
	return err
}

func setInitGlobalVariable(args []string) error {
	// set global workspace
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	dirname := "."
	if len(args) > 0 {
		dirname = args[0]
	}
	*global.WORKSPACE = filepath.Join(wd, dirname)

	return nil
}

func getProvidersList() ([]string, error) {
	var prov []string
	ui.PrintInfoLn("Getting provider list...")
	req, _ := http.NewRequest("GET", "https://github.com/selefra/registry/file-list/main/provider", nil)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return nil, err
	}
	d, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		ui.PrintErrorF("Error: %s", err.Error())
		return nil, err
	}
	d.Find(".js-navigation-open.Link--primary").Each(func(i int, s *goquery.Selection) {
		if s.Text() != "template" {
			prov = append(prov, s.Text())
		}
	})
	return prov, nil
}

func getProjectName() (projectName string) {
	defer func() {
		if projectName == "" {
			projectName = filepath.Base(*global.WORKSPACE)
		}
	}()
	var err error
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("project name:(%s)", filepath.Base(*global.WORKSPACE))

	projectName, err = reader.ReadString('\n')
	if err != nil {
		return ""
	}
	projectName = strings.TrimSpace(strings.Replace(projectName, "\n", "", -1))

	return projectName
}
