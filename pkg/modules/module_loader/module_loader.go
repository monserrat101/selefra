package module_loader

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module"
)

// ------------------------------------------------- --------------------------------------------------------------------

const (
	DownloadModulesDirectoryName = "modules"
)

// ------------------------------------------------- --------------------------------------------------------------------

// ModuleLoaderOptions Options when loading the module
type ModuleLoaderOptions struct {

	// Where can I download the module
	Source string `json:"source" yaml:"source"`

	// Which version of which module to download
	Version string `json:"version" yaml:"version"`

	// What is the download path configured in the current system
	DownloadDirectory string

	// TODO Can be used to track download progress
	ProgressTracker getter.ProgressTracker

	// It's used to send information back in real time
	MessageChannel *message.Channel[*schema.Diagnostics] `json:"message-channel"`

	// How do I go from the root module to the current module
	DependenciesTree []string
}

func (x *ModuleLoaderOptions) Copy() *ModuleLoaderOptions {
	return &ModuleLoaderOptions{
		Source:  x.Source,
		Version: x.Version,
		// TODO
		//ProgressTracker:   x.ProgressTracker,
		DownloadDirectory: x.DownloadDirectory,
		MessageChannel:    x.MessageChannel.MakeChildChannel(),
		DependenciesTree:  append([]string{}, x.DependenciesTree...),
	}
}

// ------------------------------------------------- --------------------------------------------------------------------

// ModuleLoader Module loader
type ModuleLoader[Options any] interface {

	// Name The name of the loader
	Name() ModuleLoaderType

	// Load Use this loader to load the module
	Load(ctx context.Context) (*module.Module, bool)

	Options() Options
}

// ------------------------------------------------- --------------------------------------------------------------------
