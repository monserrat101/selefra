package executors

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/env"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/modules/module_loader"
	"github.com/selefra/selefra/pkg/modules/planner"
	"github.com/selefra/selefra/pkg/providers/local_providers_manager"
	"github.com/selefra/selefra/pkg/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProviderFetchExecutor_Execute(t *testing.T) {

	projectWorkspace := "./test_data/test_fetch_module"
	downloadWorkspace := "./test_download"

	// Load the module used for the test
	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	loader, err := module_loader.NewLocalDirectoryModuleLoader(&module_loader.LocalDirectoryModuleLoaderOptions{
		ModuleLoaderOptions: &module_loader.ModuleLoaderOptions{
			MessageChannel: messageChannel,
		},
		ModuleDirectory: projectWorkspace,
	})
	assert.Nil(t, err)
	rootModule, b := loader.Load(context.Background())
	messageChannel.ReceiverWait()
	assert.NotNil(t, rootModule)
	assert.True(t, b)

	// Make an installation plan
	providersInstallPlan, diagnostics := planner.MakeProviderInstallPlan(context.Background(), rootModule)
	if utils.IsNotEmpty(diagnostics) {
		t.Log(diagnostics.ToString())
	}
	assert.False(t, utils.HasError(diagnostics))
	assert.NotNil(t, providersInstallPlan)

	// Installation-dependent dependency
	messageChannel = message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	executor, diagnostics := NewProviderInstallExecutor(&ProviderInstallExecutorOptions{
		Plans:             providersInstallPlan,
		MessageChannel:    messageChannel,
		DownloadWorkspace: downloadWorkspace,
	})
	assert.False(t, utils.HasError(diagnostics))
	if utils.IsNotEmpty(diagnostics) {
		t.Log(diagnostics.ToString())
	}
	d := executor.Execute(context.Background())
	messageChannel.ReceiverWait()
	assert.False(t, utils.HasError(d))
	if utils.IsNotEmpty(d) {
		t.Log(d.ToString())
	}

	// Develop a data pull plan
	providerFetchPlans, d := planner.NewProviderFetchPlanner(rootModule, providersInstallPlan.ToMap()).MakePlan(context.Background())
	if utils.IsNotEmpty(d) {
		t.Log(d.ToString())
	}
	assert.False(t, utils.HasError(d))
	assert.NotNil(t, providerFetchPlans)

	// Ready to start pulling
	localProviderManager, err := local_providers_manager.NewLocalProvidersManager("./test_download")
	assert.Nil(t, err)
	assert.NotNil(t, localProviderManager)
	messageChannel = message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			t.Log(message.ToString())
		}
	})
	fetchExecutor := NewProviderFetchExecutor(&ProviderFetchExecutorOptions{
		LocalProviderManager: localProviderManager,
		Plans:                providerFetchPlans,
		MessageChannel:       messageChannel,
		WorkerNum:            3,
		Workspace:            projectWorkspace,
		DSN:                  env.GetDatabaseDsn(),
	})
	d = fetchExecutor.Execute(context.Background())
	assert.False(t, utils.HasError(d))
	if utils.IsNotEmpty(d) {
		t.Log(d.ToString())
	}
	messageChannel.ReceiverWait()
	t.Log("message received")
}
