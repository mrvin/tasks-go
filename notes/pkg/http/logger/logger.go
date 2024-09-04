package logger

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type requestID string

const RequestIDKey = requestID("requestID")

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	totalWritByte int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(slByte []byte) (int, error) {
	writeByte, err := lrw.ResponseWriter.Write(slByte)
	lrw.totalWritByte += writeByte

	return writeByte, err //nolint:wrapcheck
}

func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		slog.Warn("GetRequestID: ctx is nil")
		return ""
	}
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	slog.Warn("GetRequestID: no request id in ctx")

	return ""
}

type Logger struct {
	Inner http.Handler
}

func (l *Logger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	requestID := uuid.New().String()
	ctx := context.WithValue(req.Context(), RequestIDKey, requestID)

	logReq := slog.With(
		slog.String("request_id", requestID),
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("addr", req.RemoteAddr),
	)
	timeStart := time.Now()
	lrw := NewLoggingResponseWriter(res)
	defer func() {
		logReq.Info("Request "+req.Proto,
			slog.Int("status", lrw.statusCode),
			slog.Int("bytes", lrw.totalWritByte),
			slog.String("duration", time.Since(timeStart).String()),
		)
	}()

	l.Inner.ServeHTTP(lrw, req.WithContext(ctx))
}
