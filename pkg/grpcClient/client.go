package grpcClient

import (
	"context"
	"fmt"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient/proto/issue"
	logPb "github.com/selefra/selefra/pkg/grpcClient/proto/log"
	"github.com/selefra/selefra/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strings"
	"sync"
)

type RpcClient struct {
	ctx context.Context

	conn              *grpc.ClientConn
	issueStreamClient issue.Issue_UploadIssueStreamClient
	logStreamClient   logPb.Log_UploadLogStreamClient

	taskId    string
	token     string
	statusMap map[string]string
}

var (
	client     *RpcClient
	clientOnce *sync.Once
)

func ShouldClient(ctx context.Context, token, taskId string) (*RpcClient, error) {
	if client != nil {
		return client, nil
	}

	return newClient(ctx, token, taskId)
}

func Client() *RpcClient {
	return client
}

func newClient(ctx context.Context, token, taskId string) (*RpcClient, error) {
	var err error
	clientOnce.Do(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(getDial(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return
		}

		innerClient := RpcClient{
			ctx:       ctx,
			conn:      conn,
			taskId:    taskId,
			token:     token,
			statusMap: make(map[string]string),
		}

		var openedLogStreamClient logPb.Log_UploadLogStreamClient
		logStreamClient := logPb.NewLogClient(conn)
		openedLogStreamClient, err = logStreamClient.UploadLogStream(ctx)
		if err != nil {
			return
		}
		innerClient.logStreamClient = openedLogStreamClient

		var openedIssueStreamClient issue.Issue_UploadIssueStreamClient
		issueStreamClient := issue.NewIssueClient(conn)
		openedIssueStreamClient, err = issueStreamClient.UploadIssueStream(ctx)
		if err != nil {
			return
		}
		innerClient.issueStreamClient = openedIssueStreamClient

		client = &innerClient

		utils.MultiRegisterClose(map[string]func(){
			"grpc conn": func() {
				_ = conn.Close()
			},
			"log stream": func() {
				_ = openedLogStreamClient.CloseSend()
			},
			"issue stream": func() {
				_ = openedIssueStreamClient.CloseSend()
			},
		})

		return
	})

	if err != nil {
		client = &RpcClient{}
		return nil, err
	}

	return client, nil
}

func getDial() string {
	var dialMap = make(map[string]string)
	dialMap["dev-api.selefra.io"] = "dev-tcp.selefra.io:1234"
	dialMap["main-api.selefra.io"] = "main-tcp.selefra.io:1234"
	dialMap["pre-api.selefra.io"] = "pre-tcp.selefra.io:1234"
	if dialMap[global.SERVER] != "" {
		return dialMap[global.SERVER]
	}
	arr := strings.Split(global.SERVER, ":")
	return arr[0] + ":1234"
}

func (client *RpcClient) GetIssueUploadIssueStreamClient() issue.Issue_UploadIssueStreamClient {
	return client.issueStreamClient
}

func (client *RpcClient) GetLogUploadLogStreamClient() logPb.Log_UploadLogStreamClient {
	return client.logStreamClient
}

func (client *RpcClient) SetStatus(status string) {
	if client.statusMap[global.Stage()] == "" {
		client.statusMap[global.Stage()] = status
	}
}

func (client *RpcClient) getStatus() string {
	if client.statusMap[global.Stage()] != "" {
		return client.statusMap[global.Stage()]
	}
	return "success"
}

func (client *RpcClient) newLogClient() error {
	logCli := logPb.NewLogClient(client.conn)
	uploadStreamCli, err := logCli.UploadLogStream(client.ctx)
	if err != nil {
		return err
	}
	client.logStreamClient = uploadStreamCli
	return nil
}

func (client *RpcClient) newIssueClient() error {
	issueCli := issue.NewIssueClient(client.conn)
	uploadIssueCli, err := issueCli.UploadIssueStream(client.ctx)
	if err != nil {
		return err
	}
	client.issueStreamClient = uploadIssueCli
	return nil
}

func (client *RpcClient) IssueStreamClient() issue.Issue_UploadIssueStreamClient {
	return client.issueStreamClient
}

func (client *RpcClient) LogStreamClient() logPb.Log_UploadLogStreamClient {
	return client.logStreamClient
}

func (client *RpcClient) GetTaskID() string {
	return client.taskId
}

func (client *RpcClient) GetToken() string {
	return client.token
}

func (client *RpcClient) GetConn() *grpc.ClientConn {
	return client.conn
}

func (client *RpcClient) UploadLogStatus() (*logPb.Res, error) {
	if client.conn == nil {
		return nil, nil
	}
	logClient := logPb.NewLogClient(client.conn)
	statusInfo := &logPb.StatusInfo{
		BaseInfo: &logPb.BaseConnectionInfo{
			Token:  client.GetToken(),
			TaskId: client.GetTaskID(),
		},
		Stag:   global.Stage(),
		Status: client.getStatus(),
		Time:   timestamppb.Now(),
	}
	res, err := logClient.UploadLogStatus(client.ctx, statusInfo)
	if err != nil {
		return nil, fmt.Errorf("Fail to upload log status:%s", err.Error())
	}
	return res, nil
}
