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

	ProgressTracker getter.ProgressTracker

	// It's used to send information back in real time
	MessageChannel *message.Channel[*schema.Diagnostics] `json:"message-channel"`
}

func (x *ModuleLoaderOptions) Copy() *ModuleLoaderOptions {
	return &ModuleLoaderOptions{
		Source:            x.Source,
		Version:           x.Version,
		DownloadDirectory: x.DownloadDirectory,
		MessageChannel:    x.MessageChannel.MakeChildChannel(),
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
