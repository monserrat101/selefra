package logger

import (
	"github.com/selefra/selefra/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

type Logger struct {
	logger *zap.Logger
}

func Default() *Logger {
	return defaultLogger
}

func DebugF(msg string, args ...any) {
	defaultLogger.Debug(msg, args)
}

func InfoF(msg string, args ...any) {
	defaultLogger.Info(msg, args)
}

func ErrorF(msg string, args ...any) {
	defaultLogger.Error(msg, args)
}

func FatalF(msg string, args ...any) {
	defaultLogger.Fatal(msg, args)
}

var defaultLogger, _ = NewLogger(Config{
	FileLogEnabled:    true,
	ConsoleLogEnabled: false,
	EncodeLogsAsJson:  true,
	ConsoleNoColor:    true,
	Source:            "client",
	Directory:         "logs",
	Level:             "info",
})

func NewLogger(c Config) (*Logger, error) {
	logDir := filepath.Join(global.WorkSpace(), c.Directory)
	_, err := os.Stat(logDir)
	if os.IsNotExist(err) {
		err = os.Mkdir(logDir, 0755)
	}
	if err != nil {
		return nil, nil
	}
	errorStack := zap.AddStacktrace(zap.ErrorLevel)

	development := zap.Development()

	logger := zap.New(zapcore.NewTee(c.GetEncoderCore()...), errorStack, development)

	if c.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return &Logger{logger: logger}, nil
}
