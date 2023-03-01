package main

import (
	"context"
	"fmt"
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra-utils/pkg/id_util"
	selefraGrpc "github.com/selefra/selefra/pkg/grpc"
	"github.com/selefra/selefra/pkg/grpc/pb/log"
	"github.com/selefra/selefra/pkg/message"
	"github.com/selefra/selefra/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"time"
)

func main() {

	var (
		serverUrl = os.Getenv("serverUrl")
		taskId    = os.Getenv("taskId")
		token     = os.Getenv("token")
	)
	fmt.Println("serverUrl: " + serverUrl)
	fmt.Println("taskId: " + taskId)
	fmt.Println("token: " + token)

	conn, err := grpc.Dial(serverUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	metadataContext := metadata.AppendToOutgoingContext(context.Background(), "taskUUID", taskId, "token", token)

	// create upload
	client := log.NewLogClient(conn)
	stream, err := client.UploadLogStream(metadataContext)
	if err != nil {
		panic(err)
	}

	messageChannel := message.NewChannel[*schema.Diagnostics](func(index int, message *schema.Diagnostics) {
		if utils.IsNotEmpty(message) {
			fmt.Println(message.ToString())
		}
	})
	uploaderOptions := &selefraGrpc.StreamUploaderOptions[log.Log_UploadLogStreamClient, int, *log.UploadLogStream_Request, *log.UploadLogStream_Response]{
		Name:                      "test-log-stream",
		Client:                    stream,
		WaitSendTaskQueueBuffSize: 1,
		MessageChannel:            messageChannel,
	}
	uploader := selefraGrpc.NewStreamUploader[log.Log_UploadLogStreamClient, int, *log.UploadLogStream_Request, *log.UploadLogStream_Response](uploaderOptions)
	uploader.RunUploaderWorker()

	index := 0
	for index < 30 {

		index++
		submit, diagnostics := uploader.Submit(context.Background(), index, &log.UploadLogStream_Request{
			Stage: log.StageType_STAGE_TYPE_INFRASTRUCTURE_ANALYSIS,
			Index: uint64(index),
			Msg:   fmt.Sprintf("log message body: %s", id_util.RandomId()),
			Level: log.Level_LEVEL_INFO,
			Time:  timestamppb.Now(),
		})
		if utils.IsNotEmpty(diagnostics) {
			fmt.Println(diagnostics.ToString())
		}
		fmt.Println(fmt.Sprintf("index %d, submit: %v", index, submit))

		time.Sleep(time.Second * 1)

	}

	time.Sleep(time.Second * 60)

	for index < 100 {

		index++
		submit, diagnostics := uploader.Submit(context.Background(), index, &log.UploadLogStream_Request{
			Stage: log.StageType_STAGE_TYPE_INFRASTRUCTURE_ANALYSIS,
			Index: uint64(index),
			Msg:   fmt.Sprintf("log message body: %s", id_util.RandomId()),
			Level: log.Level_LEVEL_INFO,
			Time:  timestamppb.Now(),
		})
		if utils.IsNotEmpty(diagnostics) {
			fmt.Println(diagnostics.ToString())
		}
		fmt.Println(fmt.Sprintf("index %d, submit: %v", index, submit))

		time.Sleep(time.Second * 1)

	}

	uploader.ShutdownAndWait(context.Background())
	messageChannel.ReceiverWait()

}
