package oci

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

// 在Windows平台测试，无法通过测试
func TestPostgreSQLInstaller_Run1(t *testing.T) {
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	downloader := NewPostgreSQLDownloader(&PostgreSQLDownloaderOptions{
		MessageChannel:    messageChannel,
		DownloadDirectory: "./test_download",
	})
	isRunSuccess := downloader.Run(context.Background())
	messageChannel.ReceiverWait()
	assert.True(t, isRunSuccess)
}
