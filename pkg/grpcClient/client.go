package grpcClient

import (
	"context"
	"fmt"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/grpcClient/proto/issue"
	logPb "github.com/selefra/selefra/pkg/grpcClient/proto/log"
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

		return
	})

	if err != nil {
		client = &RpcClient{}
		return nil, err
	}

	return client, nil
}

var Cli RpcClient

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

func (g *RpcClient) getDial() string {
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

func (g *RpcClient) NewConn(token, taskId string) error {
	conn, err := grpc.Dial(g.getDial(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("fail to dial: %v", err)
	}
	g.conn = conn
	g.taskId = taskId
	g.token = token
	g.statusMap = make(map[string]string)
	g.ctx = context.Background()
	err = g.newLogClient()
	if err != nil {
		return fmt.Errorf("fail to create uploadLogStreamCli cli:%s", err.Error())
	}
	err = g.newIssueClient()
	if err != nil {
		return fmt.Errorf("fail to create uploadIssueCli cli:%s", err.Error())
	}
	return err
}

func (g *RpcClient) GetIssueUploadIssueStreamClient() issue.Issue_UploadIssueStreamClient {
	return g.issueStreamClient
}

func (g *RpcClient) GetLogUploadLogStreamClient() logPb.Log_UploadLogStreamClient {
	return g.logStreamClient
}

func (client *RpcClient) SetStatus(status string) {
	if client.statusMap[global.STAG] == "" {
		client.statusMap[global.STAG] = status
	}
}

func (client *RpcClient) getStatus() string {
	if client.statusMap[global.Stage()] != "" {
		return client.statusMap[global.STAG]
	}
	return "success"
}

func (g *RpcClient) newLogClient() error {
	logCli := logPb.NewLogClient(g.conn)
	uploadStreamCli, err := logCli.UploadLogStream(g.ctx)
	if err != nil {
		return err
	}
	g.logStreamClient = uploadStreamCli
	return nil
}

func (g *RpcClient) newIssueClient() error {
	issueCli := issue.NewIssueClient(g.conn)
	uploadIssueCli, err := issueCli.UploadIssueStream(g.ctx)
	if err != nil {
		return err
	}
	g.issueStreamClient = uploadIssueCli
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
