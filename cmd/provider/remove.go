package provider

import (
	"encoding/json"
	"errors"
	"github.com/selefra/selefra/config"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/selefra/selefra/ui"
	"github.com/spf13/cobra"
	"os"
)

func newCmdProviderRemove() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "remove",
		Short:            "Remove one or more plugins from the download cache",
		Long:             "Remove one or more plugins from the download cache",
		PersistentPreRun: global.DefaultWrappedInit(),
		RunE: func(cmd *cobra.Command, names []string) error {
			return Remove(names)
		},
	}

	cmd.SetHelpFunc(cmd.HelpFunc())
	return cmd
}

func Remove(names []string) error {
	argsMap := make(map[string]bool)
	for i := range names {
		argsMap[names[i]] = true
	}
	deletedMap := make(map[string]bool)
	cof, err := config.GetConfig()
	if err != nil {
		return err
	}
	namespace, _, err := utils.Home()
	if err != nil {
		return err
	}
	provider := registry.NewProviderRegistry(namespace)

	for _, p := range cof.Selefra.Providers {
		name := *p.Source
		path := utils.GetPathBySource(*p.Source, p.Version)
		prov := registry.ProviderBinary{
			Provider: registry.Provider{
				Name:    name,
				Version: p.Version,
				Source:  "",
			},
			Filepath: path,
		}
		if !argsMap[p.Name] || deletedMap[p.Path] {
			break
		}

		err := provider.DeleteProvider(prov)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				ui.Warningf("Failed to remove  %s: %s", p.Name, err.Error())
			}
		}
		_, jsonPath, err := utils.Home()
		if err != nil {
			return err
		}
		c, err := os.ReadFile(jsonPath)
		if err == nil {
			var configMap = make(map[string]string)
			err = json.Unmarshal(c, &configMap)
			if err != nil {
				return err
			}
			delete(configMap, *p.Source+"@"+p.Version)
			c, err = json.Marshal(configMap)
			if err != nil {
				return err
			}
			err = os.Remove(jsonPath)
			if err != nil {
				return err
			}
			err = os.WriteFile(jsonPath, c, 0644)
			if err != nil {
				return err
			}
			deletedMap[path] = true
		}
		ui.Successf("Removed %s success", *p.Source)
	}
	return nil
}
