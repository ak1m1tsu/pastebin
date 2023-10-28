package logger

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/romankravchuk/pastebin/internal/controller/http/response"
	"github.com/romankravchuk/pastebin/pkg/log"
)

func New(l *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var (
				ww    = middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				start = time.Now()
			)

			defer func() {
				if rec := recover(); rec != nil {
					l.Panic(
						"recoved system error",
						log.FF{
							{Key: "recover", Value: rec},
							{Key: "duration", Value: time.Since(start)},
							{Key: "stack", Value: string(debug.Stack())},
						},
					)

					response.InternalServerError(w, r)

					return
				}

				l.Info(
					"incoming request",
					log.FF{
						{Key: "remote_addr", Value: r.RemoteAddr},
						{Key: "url", Value: r.URL.String()},
						{Key: "proto", Value: r.Proto},
						{Key: "method", Value: r.Method},
						{Key: "status", Value: ww.Status()},
						{Key: "user_agent", Value: r.UserAgent()},
						{Key: "latency", Value: time.Since(start)},
						{Key: "bytes_received", Value: r.Header.Get("Content-Length")},
						{Key: "bytes_written", Value: ww.BytesWritten()},
					},
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
