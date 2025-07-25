// Copyright (c) 2025 coze-dev Authors
// SPDX-License-Identifier: Apache-2.0

package logs

import (
	"context"
)

type LogIDProvider interface {
	GetLogID(ctx context.Context) string
	SetLogID(ctx context.Context, logID string) context.Context
	NewLogID() string
}

// Logger Interface for logging
type Logger interface {
	LogIDProvider

	SetLevel(level LogLevel)
	GetLevel() LogLevel

	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})

	CtxDebug(ctx context.Context, format string, v ...interface{})
	CtxInfo(ctx context.Context, format string, v ...interface{})
	CtxWarn(ctx context.Context, format string, v ...interface{})
	CtxError(ctx context.Context, format string, v ...interface{})
	CtxFatal(ctx context.Context, format string, v ...interface{})

	Flush()
}

// LogLevel Log Level
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)
