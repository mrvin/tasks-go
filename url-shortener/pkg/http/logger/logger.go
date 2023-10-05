package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	totalWritByte int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(slByte []byte) (writeByte int, err error) {
	writeByte, err = lrw.ResponseWriter.Write(slByte)
	lrw.totalWritByte += writeByte
	return
}

const RequestIDKey = "requestID"

type Logger struct {
	reqid uint64
	Inner http.Handler
}

func (l *Logger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	id := atomic.AddUint64(&l.reqid, 1)
	requestID := fmt.Sprintf("%06d", id)
	ctx := context.WithValue(req.Context(), RequestIDKey, requestID)

	logReq := slog.With(
		slog.String("request_id", requestID),
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("addr", req.RemoteAddr),
		//slog.String("user_agent", req.UserAgent()),
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
