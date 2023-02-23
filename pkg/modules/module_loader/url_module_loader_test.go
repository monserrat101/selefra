package module_loader

import (
	"context"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/pkg/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestURLModuleLoader_Load(t *testing.T) {
	loader, err := NewURLModuleLoader(&URLModuleLoaderOptions{
		ModuleLoaderOptions: &ModuleLoaderOptions{
			DownloadDirectory: "./test_download",
			//ProgressTracker:   testProgressTracker{},
			MessageChannel: message.NewChannel[*schema.Diagnostics](func(index int, d *schema.Diagnostics) {
				t.Log(d)
			}),
		},
		ModuleURL: "https://github.com/selefra/rules-aws-misconfiguration-s3/releases/download/v0.0.1/rules-aws-misconfigure-s3.zip",
	})
	assert.Nil(t, err)
	rootModule, b := loader.Load(context.Background())
	assert.True(t, b)
	assert.NotNil(t, rootModule)

}

//type testProgressTracker struct {
//}
//
//var _ getter.ProgressTracker = testProgressTracker{}
//
//func (x testProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) (body io.ReadCloser) {
//	fmt.Println(float64(currentSize) * 100 / float64(totalSize))
//	return stream
//}
