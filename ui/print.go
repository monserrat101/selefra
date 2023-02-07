package ui

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra/pkg/grpcClient"
	logPb "github.com/selefra/selefra/pkg/grpcClient/proto/log"
	"github.com/selefra/selefra/pkg/utils"
	"google.golang.org/protobuf/types/known/timestamppb"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	hclog "github.com/hashicorp/go-hclog"

	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"github.com/selefra/selefra/global"
	"github.com/selefra/selefra/pkg/logger"
)

type uiPrinter struct {
	// log record logs
	log *logger.Logger

	// fw is a file operator pointer for backend log file
	fw *os.File

	// rpcClient is a grpc client, it send logs to grpc server
	rpcClient *grpcClient.RpcClient

	// step store the steps for uiPrinter
	step int32
}

func newUiPrinter() *uiPrinter {
	ua := &uiPrinter{
		step: 0,
	}

	ua.log, _ = logger.NewLogger(logger.Config{
		FileLogEnabled:    true,
		ConsoleLogEnabled: false,
		EncodeLogsAsJson:  true,
		ConsoleNoColor:    true,
		Source:            "client",
		Directory:         "logs",
		Level:             "info",
	})

	flag := strings.ToLower(os.Getenv("SELEFRA_CLOUD_FLAG"))
	if flag == "true" || flag == "enable" {
		_, err := os.Stat("ws.log")
		if err != nil {
			if !os.IsNotExist(err) {
				panic("Unknown error," + err.Error())
			}
			ua.fw, err = os.Create("ws.log")
		} else {
			ua.fw, err = os.OpenFile("ws.log", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		}
		if err != nil {
			panic("ws log file open error," + err.Error())
		}
		utils.RegisterClose("ws.log", func() {
			_ = ua.fw.Close()
		})
	}

	ua.rpcClient = grpcClient.Client()

	return ua
}

var (
	printerOnce sync.Once
	printer     *uiPrinter
)

// fsync write msg to p.fw
func (p *uiPrinter) fsync(color *color.Color, msg string) {
	jsonLog := LogJSON{
		Cmd:   global.Cmd(),
		Stag:  global.Stage(),
		Msg:   msg,
		Time:  time.Now(),
		Level: getLevel(color),
	}
	byteLog, err := json.Marshal(jsonLog)
	if err != nil {
		p.log.Error(err.Error())
		return
	}

	strLog := string(byteLog)
	_, _ = p.fw.WriteString(strLog + "\n")
}

// sync do 2 things: 1. store msg to log file; 2. send msg to rpc server if rpc client exist
// sync do not show anything
func (p *uiPrinter) sync(color *color.Color, msg string) {
	// write to file
	p.fsync(color, msg)

	// send to rpc
	if p.rpcClient != nil {
		logStreamClient := p.rpcClient.LogStreamClient()
		p.step++
		if color == ErrorColor {
			p.rpcClient.SetStatus("error")
		}

		if err := logStreamClient.Send(&logPb.ConnectMsg{
			ActionName: "",
			Data: &logPb.LogJOSN{
				Cmd:   global.Cmd(),
				Stag:  global.Stage(),
				Msg:   msg,
				Time:  timestamppb.Now(),
				Level: getLevel(color),
			},
			Index: p.step,
			Msg:   "",
			BaseInfo: &logPb.BaseConnectionInfo{
				Token:  p.rpcClient.GetToken(),
				TaskId: p.rpcClient.GetTaskID(),
			},
		}); err != nil {
			p.fsync(ErrorColor, err.Error())
			return
		}
	}

	return
}

// printf The behavior of printf is like fmt.Printf that it will format the info
// when withLn is true, it will show format info with a "\n" and call sync, else without a "\n"
func (p *uiPrinter) printf(color *color.Color, format string, args ...any) {
	// logger to terminal
	if p.log != nil {
		if color == ErrorColor {
			if _, f, l, ok := runtime.Caller(2); ok {
				printer.log.Log(hclog.Error, "%s %s:%d", fmt.Sprintf(format, args...), f, l)
			}
		}
		p.log.Log(color2level(color), format, args...)
	}

	msg := fmt.Sprintf(format, args...)

	p.sync(color, msg)

	_, _ = color.Printf(format, args...)

}

// println The behavior of println is like fmt.Println
// it will show the log info and then call sync
func (p *uiPrinter) println(color *color.Color, args ...any) {
	// logger to terminal
	if p.log != nil {
		if color == ErrorColor {
			if _, f, l, ok := runtime.Caller(2); ok {
				printer.log.Log(hclog.Error, "%s %s:%d", fmt.Sprintln(args...), f, l)
			}
		}
		p.log.Log(color2level(color), fmt.Sprintln(args...))
	}

	msg := fmt.Sprint(args...)

	p.sync(color, msg)

	_, _ = color.Println(args...)

	return
}

func color2level(color *color.Color) hclog.Level {
	switch color {
	case ErrorColor:
		return hclog.Error
	case WarningColor:
		return hclog.Warn
	case InfoColor:
		return hclog.Info
	case SuccessColor:
		return hclog.Info
	default:
		return hclog.Info
	}
}

var levelMap = map[string]int{
	"trace":   0,
	"debug":   1,
	"info":    2,
	"warning": 3,
	"error":   4,
	"fatal":   5,
}

var levelColor = []*color.Color{
	InfoColor,
	InfoColor,
	InfoColor,
	WarningColor,
	ErrorColor,
	ErrorColor,
}

// var step int32 = 0
var defaultLogger, _ = logger.NewLogger(logger.Config{
	FileLogEnabled:    true,
	ConsoleLogEnabled: false,
	EncodeLogsAsJson:  true,
	ConsoleNoColor:    true,
	Source:            "client",
	Directory:         "logs",
	Level:             "info",
})

func StoLogger() (*logger.StoLogger, error) {
	return logger.NewStoLogger(logger.Config{
		FileLogEnabled:    true,
		ConsoleLogEnabled: false,
		EncodeLogsAsJson:  true,
		ConsoleNoColor:    true,
		Source:            "client",
		Directory:         "logs",
		Level:             "info",
	})
}

func init() {
	printerOnce.Do(func() {
		printer = newUiPrinter()
	})
}

const (
	prefixManaged   = "managed"
	prefixUnmanaged = "unmanaged"
	defaultAlias    = "default"
)

var (
	ErrorColor   = color.New(color.FgRed, color.Bold)
	WarningColor = color.New(color.FgYellow, color.Bold)
	InfoColor    = color.New(color.FgWhite, color.Bold)
	SuccessColor = color.New(color.FgGreen, color.Bold)
)

type LogJSON struct {
	Cmd   string    `json:"cmd"`
	Stag  string    `json:"stag"`
	Msg   string    `json:"msg"`
	Time  time.Time `json:"time"`
	Level string    `json:"level"`
}

func getLevel(c *color.Color) string {
	var level string
	switch c {
	case ErrorColor:
		level = "error"
	case WarningColor:
		level = "warn"
	case InfoColor:
		level = "info"
	case SuccessColor:
		level = "success"
	default:
	}
	return level
}

func Errorf(format string, a ...interface{}) {
	printer.printf(ErrorColor, format, a...)
}

func Warningf(format string, a ...interface{}) {
	printer.printf(WarningColor, format, a...)
}

func Successf(format string, a ...interface{}) {
	printer.printf(SuccessColor, format, a...)
}

func Infof(format string, a ...interface{}) {
	printer.printf(InfoColor, format, a...)
}

func Errorln(a ...interface{}) {
	printer.println(ErrorColor, a...)
}

func Warningln(a ...interface{}) {
	printer.println(WarningColor, a...)
}

func Successln(a ...interface{}) {
	printer.println(SuccessColor, a...)
}

func Infoln(a ...interface{}) {
	printer.println(InfoColor, a...)
}

func Printf(c *color.Color, format string, a ...any) {
	printer.printf(c, format, a...)
}

func Println(c *color.Color, a ...any) {
	printer.println(c, a...)
}

func Print(msg string, show bool) {
	if show {
		Infoln(msg)
		return
	}

	printer.sync(InfoColor, msg)
}

func SaveLogToDiagnostic(diagnostics []*schema.Diagnostic) {
	for i := range diagnostics {
		if int(diagnostics[i].Level()) >= levelMap[global.LOGLEVEL] {
			defaultLogger.Log(hclog.Level(levelMap[global.LOGLEVEL]+1), diagnostics[i].Content())
		}
	}
}

func PrintDiagnostic(diagnostics []*schema.Diagnostic) error {
	var err error
	for i := range diagnostics {
		if int(diagnostics[i].Level()) >= levelMap[global.LOGLEVEL] {
			defaultLogger.Log(hclog.Level(levelMap[global.LOGLEVEL]+1), diagnostics[i].Content())
			Println(levelColor[int(diagnostics[i].Level())], diagnostics[i].Content())
			if diagnostics[i].Level() == schema.DiagnosisLevelError {
				err = errors.New(diagnostics[i].Content())
			}
		}
	}
	return err
}
