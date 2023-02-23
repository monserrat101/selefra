package apply

//import (
//	"github.com/selefra/selefra/config"
//	"github.com/selefra/selefra/global"
//	"github.com/stretchr/testify/require"
//	"testing"
//)
//
////func TestApply(t *testing.T) {
////	global.Init("TestApply", global.WithWorkspace("../../tests/workspace/offline"))
////	err := Apply(context.Background())
////	if err != nil {
////		t.Error(err)
////	}
////}
//
//func Test_RunPathModule(t *testing.T) {
//	global.Init("", global.WithWorkspace("../../tests/workspace/offline"))
//	modules, err := config.GetModules()
//	require.NoError(t, err, "get modules failed")
//
//	require.Equal(t, "Misconfigure-S3", modules[0].Name)
//
//	rules := GetModuleRules(modules[0])
//
//	require.Equal(t, 1, len(rules), "rules length error")
//	require.Equal(t, "ebs_volume_are_unencrypted", rules[0].Name, "rules name error")
//}
//
//func Test_RunRulesWithoutModule(t *testing.T) {
//	global.Init("", global.WithWorkspace("../../tests/workspace/offline"))
//
//	rules := GetAllRules()
//
//	require.Equal(t, 1, len(rules), "rules length error")
//	require.Equal(t, "ebs_volume_are_unencrypted", rules[0].Name, "rules name error")
//}
