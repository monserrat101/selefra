package module_loader

import (
	"context"
	"errors"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/modules/parser"
	"github.com/selefra/selefra/pkg/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ------------------------------------------------- --------------------------------------------------------------------

// LocalDirectoryModuleLoaderOptions Option when loading modules from a local directory
type LocalDirectoryModuleLoaderOptions struct {
	*ModuleLoaderOptions

	// Directory where the module resides Directory
	ModuleDirectory string `json:"module-directory" yaml:"module-directory"`
}

func (x *LocalDirectoryModuleLoaderOptions) Copy() *LocalDirectoryModuleLoaderOptions {
	return &LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: x.ModuleLoaderOptions.Copy(),
		ModuleDirectory:     x.ModuleDirectory,
	}
}

func (x *LocalDirectoryModuleLoaderOptions) CopyForModuleDirectory(moduleDirectory string) *LocalDirectoryModuleLoaderOptions {
	options := x.Copy()
	options.ModuleDirectory = moduleDirectory
	return options
}

// ------------------------------------------------- --------------------------------------------------------------------

// LocalDirectoryModuleLoader Load the module from the local directory
type LocalDirectoryModuleLoader struct {
	options *LocalDirectoryModuleLoaderOptions
}

var _ ModuleLoader[*LocalDirectoryModuleLoaderOptions] = &LocalDirectoryModuleLoader{}

func NewLocalDirectoryModuleLoader(options *LocalDirectoryModuleLoaderOptions) (*LocalDirectoryModuleLoader, error) {

	if !utils.ExistsDirectory(options.ModuleDirectory) {
		return nil, fmt.Errorf("module directory %s does not exist or is not directory", options.ModuleDirectory)
	}

	return &LocalDirectoryModuleLoader{
		options: options,
	}, nil
}

func (x *LocalDirectoryModuleLoader) Name() ModuleLoaderType {
	return ModuleLoaderTypeLocalDirectory
}

func (x *LocalDirectoryModuleLoader) Load(ctx context.Context) (*module.Module, bool) {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
	}()

	diagnostics := schema.NewDiagnostics()

	// check path
	if diagnostics.AddDiagnostics(x.checkModuleDirectory()).HasError() {
		x.options.MessageChannel.Send(diagnostics)
		return nil, false
	}

	// list all yaml file
	yamlFilePathSlice, d := x.listModuleDirectoryYamlFilePath()
	if diagnostics.AddDiagnostics(d).HasError() {
		x.options.MessageChannel.Send(diagnostics)
		return nil, false
	}

	// Read all files under the module as modules, these modules may be incomplete, may be some fragments of the module
	yamlFileModuleSlice := make([]*module.Module, len(yamlFilePathSlice))
	isHasError := false
	for index, yamlFilePath := range yamlFilePathSlice {
		yamlFileModule, d := parser.NewYamlFileToModuleParser(yamlFilePath).Parse()
		if utils.HasError(d) {
			isHasError = true
		}
		x.options.MessageChannel.Send(d)
		yamlFileModuleSlice[index] = yamlFileModule
	}
	if isHasError {
		return nil, false
	}

	// Merge these modules
	finalModule := &module.Module{}
	hasError := false
	for _, yamlFileModule := range yamlFileModuleSlice {
		merge, d := finalModule.Merge(yamlFileModule)
		if d != nil && d.HasError() {
			hasError = true
		}
		x.options.MessageChannel.Send(d)
		if merge != nil {
			finalModule = merge
		}
	}
	if hasError {
		return nil, false
	}

	// Print phased results
	x.options.MessageChannel.Send(schema.NewDiagnostics().AddInfo("load module from %s success", x.options.ModuleDirectory))

	// load sub modules
	subModuleSlice, loadSuccess := x.loadSubModules(ctx, finalModule.ModulesBlock)
	if !loadSuccess {
		return nil, false
	}
	finalModule.SubModules = subModuleSlice

	return finalModule, true
}

func (x *LocalDirectoryModuleLoader) loadSubModules(ctx context.Context, modulesBlock module.ModulesBlock) ([]*module.Module, bool) {
	subModuleSlice := make([]*module.Module, 0)
	for _, moduleBlock := range modulesBlock {
		for index, useModuleSource := range moduleBlock.Uses {

			useLocation := moduleBlock.GetNodeLocation(fmt.Sprintf("uses[%d]._value", index))
			moduleDirectoryPath := filepath.Dir(useLocation.Path)

			switch NewModuleLoaderBySource(useModuleSource) {
			case ModuleLoaderTypeInvalid:
				errorReport := module.RenderErrorTemplate(fmt.Sprintf("invalid module uses source %s, unsupported module loader", useModuleSource), useLocation)
				x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
				return nil, false
			case ModuleLoaderTypeS3Bucket:
				s3BucketModuleLoaderOptions := &S3BucketModuleLoaderOptions{
					ModuleLoaderOptions: x.options.ModuleLoaderOptions.Copy(),
					S3BucketURL:         useModuleSource,
				}
				loader, err := NewS3BucketModuleLoader(s3BucketModuleLoaderOptions)
				if err != nil {
					s3BucketModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
					errorReport := module.RenderErrorTemplate(fmt.Sprintf("create s3 module loader error: %s", err.Error()), useLocation)
					x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
					return nil, false
				}
				subModule, loadSuccess := loader.Load(ctx)
				if !loadSuccess {
					return nil, false
				}
				subModuleSlice = append(subModuleSlice, subModule)

			// TODO 2023-2-20 15:31:17 Comment this out for now. Tuning this process takes too long and may not be possible
			//case ModuleLoaderTypeGitHubRegistry:
			//	gitHubRegistryModuleLoaderOptions := &GitHubRegistryModuleLoaderOptions{
			//		ModuleLoaderOptions:  x.options.ModuleLoaderOptions.Copy(),
			//		RegistryRepoFullName: useModuleSource,
			//	}
			//	loader, err := NewGitHubRegistryModuleLoader(gitHubRegistryModuleLoaderOptions)
			//	if err != nil {
			//		gitHubRegistryModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
			//		errorReport := module.RenderErrorTemplate(fmt.Sprintf("create github registry module loader error: %s", err.Error()), useLocation)
			//		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
			//		return nil, false
			//	}
			//	subModule, loadSuccess := loader.Load(ctx)
			//	if !loadSuccess {
			//		return nil, false
			//	}
			//	subModuleSlice = append(subModuleSlice, subModule)

			case ModuleLoaderTypeLocalDirectory:
				// The path of the submodule should be from the current path
				submoduleDirectoryPath := filepath.Join(moduleDirectoryPath, useModuleSource)
				localDirectoryModuleLoaderOptions := x.options.CopyForModuleDirectory(submoduleDirectoryPath)
				loader, err := NewLocalDirectoryModuleLoader(localDirectoryModuleLoaderOptions)
				if err != nil {
					localDirectoryModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
					errorReport := module.RenderErrorTemplate(fmt.Sprintf("create local directory module loader error: %s", err.Error()), useLocation)
					x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
					return nil, false
				}
				subModule, loadSuccess := loader.Load(ctx)
				if !loadSuccess {
					return nil, false
				}
				subModuleSlice = append(subModuleSlice, subModule)

			case ModuleLoaderTypeURL:
				urlModuleLoaderOptions := &URLModuleLoaderOptions{
					ModuleLoaderOptions: x.options.ModuleLoaderOptions.Copy(),
					ModuleURL:           useModuleSource,
				}
				loader, err := NewURLModuleLoader(urlModuleLoaderOptions)
				if err != nil {
					urlModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
					errorReport := module.RenderErrorTemplate(fmt.Sprintf("create url module loader error: %s", err.Error()), useLocation)
					x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
					return nil, false
				}
				subModule, loadSuccess := loader.Load(ctx)
				if !loadSuccess {
					return nil, false
				}
				subModuleSlice = append(subModuleSlice, subModule)

			default:
				errorReport := module.RenderErrorTemplate(fmt.Sprintf("module source %s can  cannot be assign loader", useModuleSource), useLocation)
				x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg(errorReport))
				return nil, false
			}
		}
	}
	return subModuleSlice, true
}

func (x *LocalDirectoryModuleLoader) Options() *LocalDirectoryModuleLoaderOptions {
	return x.options
}

// Check that the given module local path is correct
func (x *LocalDirectoryModuleLoader) checkModuleDirectory() *schema.Diagnostics {
	info, err := os.Stat(x.options.ModuleDirectory)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return schema.NewDiagnostics().AddErrorMsg("%s module directory %s not found", x.Name(), x.options.ModuleDirectory)
		} else {
			return schema.NewDiagnostics().AddErrorMsg("%s module directory %s, error message: %s", x.Name(), x.options.ModuleDirectory, err.Error())
		}
	}

	if !info.IsDir() {
		return schema.NewDiagnostics().AddErrorMsg("%s module directory %s found, but not directory", x.Name(), x.options.ModuleDirectory)
	}

	return nil
}

// Lists all yaml files in the directory where the module resides
func (x *LocalDirectoryModuleLoader) listModuleDirectoryYamlFilePath() ([]string, *schema.Diagnostics) {
	dir, err := os.ReadDir(x.options.ModuleDirectory)
	if err != nil {
		return nil, schema.NewDiagnostics().AddErrorMsg("%s module directory %s visit error: ", x.Name(), x.options.ModuleDirectory)
	}
	yamlFileSlice := make([]string, 0)
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		if IsYamlFile(entry) {
			yamlFilePath := filepath.Join(x.options.ModuleDirectory, entry.Name())
			yamlFileSlice = append(yamlFileSlice, yamlFilePath)
		}
	}
	return yamlFileSlice, nil
}

func IsYamlFile(entry os.DirEntry) bool {
	if entry.IsDir() {
		return false
	}
	ext := strings.ToLower(path.Ext(entry.Name()))
	return strings.HasSuffix(ext, ".yaml") || strings.HasSuffix(ext, ".yml")
}
