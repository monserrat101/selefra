package module_loader

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/md5_util"
	"github.com/selefra/selefra/pkg/http_client"
	"github.com/selefra/selefra/pkg/modules/module"
	"path/filepath"
)

// ------------------------------------------------- --------------------------------------------------------------------

// URLModuleLoaderOptions Parameter options when creating the URL module loader
type URLModuleLoaderOptions struct {
	*ModuleLoaderOptions

	// Module URL, It's a zip package
	ModuleURL string
}

func (x *URLModuleLoaderOptions) Copy() *URLModuleLoaderOptions {
	return &URLModuleLoaderOptions{
		ModuleLoaderOptions: x.ModuleLoaderOptions.Copy(),
		ModuleURL:           x.ModuleURL,
	}
}

func (x *URLModuleLoaderOptions) CopyForURL(moduleURL string) *URLModuleLoaderOptions {
	options := x.Copy()
	options.ModuleURL = moduleURL
	return options
}

// ------------------------------------------------- --------------------------------------------------------------------

// URLModuleLoader 从一个URL加载模块，这个模块应该是个压缩包，解压之后正好是模块的目录
type URLModuleLoader struct {
	options *URLModuleLoaderOptions

	// 要下载到哪个路径
	moduleDownloadDirectoryPath string
}

var _ ModuleLoader[*URLModuleLoaderOptions] = &URLModuleLoader{}

func NewURLModuleLoader(options *URLModuleLoaderOptions) (*URLModuleLoader, error) {

	directoryName, err := md5_util.Md5String(options.ModuleURL)
	if err != nil {
		// TODO
		return nil, err
	}
	moduleDownloadDirectoryPath := filepath.Join(options.DownloadDirectory, DownloadModulesDirectoryName, string(ModuleLoaderTypeURL), directoryName)

	return &URLModuleLoader{
		options:                     options,
		moduleDownloadDirectoryPath: moduleDownloadDirectoryPath,
	}, nil
}

func (x *URLModuleLoader) Name() ModuleLoaderType {
	return ModuleLoaderTypeURL
}

func (x *URLModuleLoader) Load(ctx context.Context) (*module.Module, bool) {

	defer func() {
		x.options.MessageChannel.SenderWaitAndClose()
	}()

	// step 01. Download and decompress
	err := http_client.DownloadToDirectory(ctx, x.moduleDownloadDirectoryPath, x.options.ModuleURL, x.options.ProgressTracker)
	if err != nil {
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("url module loader, url = %s, error = %s", x.options.ModuleURL, err.Error()))
		return nil, false
	}

	// step 02. The download is decompressed and converted to loading from the local path
	localDirectoryModuleLoaderOptions := &LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: x.options.ModuleLoaderOptions.Copy(),
		ModuleDirectory:     x.moduleDownloadDirectoryPath,
	}
	loader, err := NewLocalDirectoryModuleLoader(localDirectoryModuleLoaderOptions)
	if err != nil {
		localDirectoryModuleLoaderOptions.MessageChannel.SenderWaitAndClose()
		x.options.MessageChannel.Send(schema.NewDiagnostics().AddErrorMsg("create local directory %s module loader error: %s", x.moduleDownloadDirectoryPath, err.Error()))
		return nil, false
	}

	return loader.Load(ctx)
}

func (x *URLModuleLoader) Options() *URLModuleLoaderOptions {
	return x.options
}
