package module_loader

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalDirectoryModuleLoader_Load(t *testing.T) {
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if message != nil && !message.IsEmpty() {
			t.Log(message.ToString())
		}
	})
	loader, err := NewLocalDirectoryModuleLoader(&LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: &ModuleLoaderOptions{
			MessageChannel:    messageChannel,
			DownloadDirectory: testDownloadDirectory,
			Source:            "rules-aws-misconfigure-s3@v0.0.1",
			Version:           "",
		},
		ModuleDirectory: "./test_data",
	})
	assert.Nil(t, err)
	rootModule, isLoadSuccess := loader.Load(context.Background())
	assert.True(t, isLoadSuccess)
	assert.NotNil(t, rootModule)
	messageChannel.ReceiverWait()
}
