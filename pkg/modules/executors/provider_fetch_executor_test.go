package executors

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderFetchExecutor_Execute(t *testing.T) {

	projectWorkspace := "./test_data/test_fetch_module"
	downloadWorkspace := "./test_download"

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	d := NewProjectLocalLifeCycleExecutor(&ProjectLocalLifeCycleExecutorOptions{
		ProjectWorkspace:                     projectWorkspace,
		DownloadWorkspace:                    downloadWorkspace,
		MessageChannel:                       messageChannel,
		ProjectLifeCycleStep:                 ProjectLifeCycleStepFetch,
		FetchStep:                            FetchStepFetch,
		ProjectCloudLifeCycleExecutorOptions: nil,
		DSN:                                  env.GetDatabaseDsn(),
		FetchWorkerNum:                       1,
		QueryWorkerNum:                       1,
	}).Execute(context.Background())
	if utils.IsNotEmpty(d) {
		t.Log(d.ToString())
	}
	assert.False(t, utils.HasError(d))

}

func TestProviderFetchExecutorWorker_computeNeedFetchTables(t *testing.T) {

	projectWorkspace := "./test_data/test_fetch_module_with_cache"
	downloadWorkspace := "./test_download"

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	d := NewProjectLocalLifeCycleExecutor(&ProjectLocalLifeCycleExecutorOptions{
		ProjectWorkspace:                     projectWorkspace,
		DownloadWorkspace:                    downloadWorkspace,
		MessageChannel:                       messageChannel,
		ProjectLifeCycleStep:                 ProjectLifeCycleStepFetch,
		FetchStep:                            FetchStepFetch,
		ProjectCloudLifeCycleExecutorOptions: nil,
		DSN:                                  env.GetDatabaseDsn(),
		FetchWorkerNum:                       1,
		QueryWorkerNum:                       1,
	}).Execute(context.Background())
	if utils.IsNotEmpty(d) {
		t.Log(d.ToString())
	}
	assert.False(t, utils.HasError(d))

}
