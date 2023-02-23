package local_providers_manager

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/utils"
)

// Get 获取本地已经安装的provider的信息
func (x *LocalProvidersManager) Get(ctx context.Context, localProvider *LocalProvider) (*LocalProvider, *schema.Diagnostics) {
	diagnostics := schema.NewDiagnostics()
	providerVersionMetaFilePath := x.buildLocalProviderVersionMetaFilePath(localProvider.Name, localProvider.Version)
	localProviderMeta, err := utils.ReadJsonFile[*LocalProvider](providerVersionMetaFilePath)
	if err != nil {
		return nil, diagnostics.AddErrorMsg("read local provider version %s meta file failed: %s", localProvider.String(), err)
	}
	return localProviderMeta, diagnostics
}
