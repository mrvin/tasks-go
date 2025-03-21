package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type contextKey int

const (
	contextKeyRequestID contextKey = iota
	contextKeyUserID
	contextKeyScheduleID
)

const logFileMode = 0755

type Conf struct {
	FilePath string
	Level    string
}

type ContextHandler struct {
	slog.Handler
}

func Init(conf *Conf) (*os.File, error) {
	var level slog.Level

	if err := level.UnmarshalText([]byte(conf.Level)); err != nil {
		return nil, fmt.Errorf("getting level from text: %w", err)
	}

	var err error
	var logFile *os.File
	if conf.FilePath == "" {
		logFile = os.Stdout
	} else {
		logFile, err = os.OpenFile(conf.FilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, logFileMode)
		if err != nil {
			return nil, fmt.Errorf("failed open log file: %w", err)
		}
	}

	replaceAttr := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Any().(time.Time) //nolint:forcetypeassert
			a.Value = slog.StringValue(t.Format(time.StampNano))
		}
		return a
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource:   true,
		Level:       level,
		ReplaceAttr: replaceAttr,
	})

	logger := slog.New(ContextHandler{handler})
	slog.SetDefault(logger)

	return logFile, nil
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, requestID)
}

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, contextKeyUserID, userID)
}

func WithScheduleID(ctx context.Context, scheduleID int64) context.Context {
	return context.WithValue(ctx, contextKeyScheduleID, scheduleID)
}

func (h ContextHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value(contextKeyRequestID).(string); ok {
		r.Add("request_id", requestID)
	}
	if userID, ok := ctx.Value(contextKeyUserID).(string); ok {
		r.Add("user_id", userID)
	}
	if scheduleID, ok := ctx.Value(contextKeyScheduleID).(int64); ok {
		r.Add("schedule_id", scheduleID)
	}

	return h.Handler.Handle(ctx, r) //nolint:wrapcheck
}
