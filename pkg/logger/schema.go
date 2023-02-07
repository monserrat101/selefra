package logger

import (
	"github.com/selefra/selefra-provider-sdk/provider/schema"
	"go.uber.org/zap"
)

// SchemaLoggerImpl is the implement of schema.ClientLogger
type SchemaLoggerImpl struct {
	l *Logger
}

var _ schema.ClientLogger = (*SchemaLoggerImpl)(nil)

func NewSchemaLoggeer() *SchemaLoggerImpl {
	return &SchemaLoggerImpl{
		l: defaultLogger,
	}
}

func (s *SchemaLoggerImpl) Debug(msg string, fields ...zap.Field) {
	s.l.Debug(msg, fields)
}

func (s *SchemaLoggerImpl) DebugF(msg string, args ...any) {
	s.l.Debug(msg, args)
}

func (s *SchemaLoggerImpl) Info(msg string, fields ...zap.Field) {
	s.l.Info(msg, fields)
}

func (s *SchemaLoggerImpl) InfoF(msg string, args ...any) {
	s.l.Info(msg, args)
}

func (s *SchemaLoggerImpl) Warn(msg string, fields ...zap.Field) {
	s.l.Warn(msg, fields)
}

func (s *SchemaLoggerImpl) WarnF(msg string, args ...any) {
	s.l.Warn(msg, args)
}

func (s *SchemaLoggerImpl) Error(msg string, fields ...zap.Field) {
	s.l.Error(msg, fields)
}

func (s *SchemaLoggerImpl) ErrorF(msg string, args ...any) {
	s.l.Error(msg, args)
}

func (s *SchemaLoggerImpl) Fatal(msg string, fields ...zap.Field) {
	s.l.Fatal(msg, fields)
}

func (s *SchemaLoggerImpl) FatalF(msg string, args ...any) {
	s.l.Fatal(msg, args)
}
