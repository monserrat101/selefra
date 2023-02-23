package module_loader

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

const testDownloadDirectory = "./test_download"

func TestGitHubRegistryModuleLoader_Load(t *testing.T) {

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		t.Log(message.ToString())
	})

	loader, err := NewGitHubRegistryModuleLoader(&GitHubRegistryModuleLoaderOptions{
		ModuleLoaderOptions: &ModuleLoaderOptions{
			MessageChannel:    messageChannel,
			DownloadDirectory: testDownloadDirectory,
			Source:            "rules-aws-misconfigure-s3@v0.0.1",
			Version:           "",
		},
		RegistryRepoFullName: "selefra/registry",
	})
	assert.Nil(t, err)
	rootModule, b := loader.Load(context.Background())
	assert.True(t, b)
	assert.NotNil(t, rootModule)

	messageChannel.ReceiverWait()

}
