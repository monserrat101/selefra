package test

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/utils"
	"testing"
)

func TestTestCommandExecutor_Run(t *testing.T) {

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})

	NewTestCommandExecutor(&TestCommandExecutorOptions{
		ProjectWorkspace:  "./test_data",
		DownloadWorkspace: "./test_download",
		MessageChannel:    messageChannel,
		DSN:               env.GetDatabaseDsn(),
	}).Run(context.Background())
	messageChannel.ReceiverWait()
}
