package apply

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApply(t *testing.T) {
	projectWorkspace := "./test_data/test_query_module"
	downloadWorkspace := "./test_download"
	err := Apply(context.Background(), projectWorkspace, downloadWorkspace)
	assert.Nil(t, err)
}
