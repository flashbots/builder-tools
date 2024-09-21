// Package common contains common utilities and functions used by the service.
package common

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

type LoggingOpts struct {
	Debug   bool
	JSON    bool
	Service string
	Version string
}

func SetupLogger(opts *LoggingOpts) (log *slog.Logger) {
	logLevel := slog.LevelInfo
	if opts.Debug {
		logLevel = slog.LevelDebug
	}

	if opts.JSON {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	} else {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	}

	if opts.Service != "" {
		log = log.With("service", opts.Service)
	}

	if opts.Version != "" {
		log = log.With("version", opts.Version)
	}

	return log
}

// LoggingMiddlewareSlog logs the incoming HTTP request & its duration.
func LoggingMiddlewareSlog(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					method := ""
					url := ""
					if r != nil {
						method = r.Method
						url = r.URL.EscapedPath()
					}

					logger.Error(fmt.Sprintf("http request panic: %s %s", method, url),
						"err", err,
						"trace", string(debug.Stack()),
						"method", r.Method,
					)
				}
			}()
			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Info(fmt.Sprintf("http: %s %s %d", r.Method, r.URL.EscapedPath(), wrapped.status),
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", fmt.Sprintf("%f", time.Since(start).Seconds()),
				"durationUs", strconv.FormatInt(time.Since(start).Microseconds(), 10),
			)
		},
	)
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}
