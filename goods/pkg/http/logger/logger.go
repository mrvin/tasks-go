package logger

import (
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode    int
	totalWritByte int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(slByte []byte) (writeByte int, err error) { //nolint:nonamedreturns
	writeByte, err = lrw.ResponseWriter.Write(slByte)
	lrw.totalWritByte += writeByte
	return
}

type Logger struct {
	Inner http.Handler
}

func (l *Logger) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	logReq := slog.With(
		slog.String("method", req.Method),
		slog.String("path", req.URL.Path),
		slog.String("addr", req.RemoteAddr),
	)
	timeStart := time.Now()
	lrw := newLoggingResponseWriter(res)
	defer func() {
		logReq.Info("Request "+req.Proto,
			slog.Int("status", lrw.statusCode),
			slog.Int("bytes", lrw.totalWritByte),
			slog.String("duration", time.Since(timeStart).String()),
		)
	}()

	l.Inner.ServeHTTP(lrw, req)
}
