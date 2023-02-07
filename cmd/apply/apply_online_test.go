package apply

//func TestApplyOnLine(t *testing.T) {
//	global.Init("TestApplyOnLine", global.WithWorkspace("../../tests/workspace/online"))
//	_ = login.ShouldLogin()
//	defer func() {
//		logCli := grpcClient.Cli.GetLogUploadLogStreamClient()
//		conn := grpcClient.Cli.GetConn()
//		if logCli != nil {
//			err := logCli.CloseSend()
//			if err != nil {
//				log.Fatalf("fail to close log stream:%s", err.Error())
//			}
//		}
//		if conn != nil {
//			err := conn.Close()
//			if err != nil {
//				log.Fatalf("fail to close grpc conn:%s", err.Error())
//			}
//		}
//	}()
//	if testing.Short() {
//		t.Skip("skipping test in short mode.")
//		return
//	}
//	err := Apply(context.Background())
//	if err != nil {
//		t.Error(err)
//	}
//}
