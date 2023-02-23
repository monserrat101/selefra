package fetch

//
//func TestFetchOnline(t *testing.T) {
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//		return
//	}
//	global.WorkSpace() = "../../tests/workspace/online"
//	global.SERVER = "dev-api.selefra.io"
//	global.LOGINTOKEN = "xxxxxxxxxxxxxxxxxxxxxx"
//	ctx := context.Background()
//	var cof = &config.SelefraConfig{}
//	err := cof.GetConfig()
//	for _, p := range cof.Selefra.RequireProvidersBlock {
//		err = Fetch(ctx, cof, p)
//		if err != nil {
//			t.Error(err)
//		}
//	}
//}
