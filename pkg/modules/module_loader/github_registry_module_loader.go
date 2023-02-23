package module_loader

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/pointer"
	"github.com/selefra/selefra/pkg/modules/module"
	"github.com/selefra/selefra/pkg/registry"
	"github.com/selefra/selefra/pkg/version"
	"path/filepath"
)

// ------------------------------------------------- --------------------------------------------------------------------

// GitHubRegistryModuleLoaderOptions Options when creating the GitHub Registry
type GitHubRegistryModuleLoaderOptions struct {
	*ModuleLoaderOptions

	// The full name of the Registry's repository
	RegistryRepoFullName string `json:"registry-repo-full-name" yaml:"registry-repo-full-name"`
}

// ------------------------------------------------- --------------------------------------------------------------------

// GitHubRegistryModuleLoader Load the module from GitHub's Registry
type GitHubRegistryModuleLoader struct {
	githubRegistry *registry.ModuleGitHubRegistry
	options        *GitHubRegistryModuleLoaderOptions

	downloadModule          *registry.Module
	moduleDownloadDirectory string
}

var _ ModuleLoader[*GitHubRegistryModuleLoaderOptions] = &GitHubRegistryModuleLoader{}

func NewGitHubRegistryModuleLoader(options *GitHubRegistryModuleLoaderOptions) (*GitHubRegistryModuleLoader, error) {

	registryOptions := registry.NewModuleGithubRegistryOptions(options.DownloadDirectory, options.RegistryRepoFullName)
	githubRegistry, err := registry.NewModuleGitHubRegistry(registryOptions)
	if err != nil {
		return nil, err
	}

	// check params
	moduleNameAndVersion := version.ParseNameAndVersion(options.Source)
	metadata, err := githubRegistry.GetMetadata(context.Background(), registry.NewModule(moduleNameAndVersion.Name, moduleNameAndVersion.Version))
	if err != nil {
		return nil, err
	}
	moduleVersion := moduleNameAndVersion.Version
	if moduleNameAndVersion.IsLatestVersion() {
		moduleVersion = metadata.LatestVersion
	}

	if !metadata.HasVersion(moduleVersion) {
		// TODO
		return nil, fmt.Errorf("module version not found ")
	}

	// The version to which the module will be downloaded
	moduleDownloadDirectory := filepath.Join(options.DownloadDirectory, registry.ModulesListDirectoryName, moduleNameAndVersion.Name, moduleVersion)

	return &GitHubRegistryModuleLoader{
		githubRegistry:          githubRegistry,
		options:                 options,
		downloadModule:          registry.NewModule(moduleNameAndVersion.Name, moduleVersion),
		moduleDownloadDirectory: moduleDownloadDirectory,
	}, nil
}

func (x *GitHubRegistryModuleLoader) Name() ModuleLoaderType {
	return ModuleLoaderTypeGitHubRegistry
}

func (x *GitHubRegistryModuleLoader) Load(ctx context.Context) (*module.Module, bool) {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
	}()

	// Download the given repository
	downloadOptions := &registry.ModuleRegistryDownloadOptions{
		ModuleDownloadDirectoryPath: x.moduleDownloadDirectory,
		SkipVerify:                  pointer.TruePointer(),
		ProgressTracker:             x.options.ProgressTracker,
	}
	moduleDownloadDirectory, err := x.githubRegistry.Download(ctx, x.downloadModule, downloadOptions)
	if err != nil {
		// TODO
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("from github registry download module %s failed: %s", x.downloadModule.String(), err.Error()))
		return nil, false
	}

	// Continue to load submodules, if any
	localDirectoryModuleLoaderOptions := &LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: x.options.ModuleLoaderOptions.Copy(),
		ModuleDirectory:     moduleDownloadDirectory,
	}
	loader, err := NewLocalDirectoryModuleLoader(localDirectoryModuleLoaderOptions)
	if err != nil {
		localDirectoryModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("new local directory %s module loader failed: %s", moduleDownloadDirectory, err.Error()))
		return nil, false
	}

	return loader.Load(ctx)
}

func (x *GitHubRegistryModuleLoader) Options() *GitHubRegistryModuleLoaderOptions {
	return x.options
}

// ------------------------------------------------- --------------------------------------------------------------------
