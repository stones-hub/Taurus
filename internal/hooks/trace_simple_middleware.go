package hooks

import (
	"Taurus/pkg/contextx"
	"Taurus/pkg/logx"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// 重写http.ResponseWriter
type traceResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// 包装http.ResponseWriter
func WrapResponseWriter(w http.ResponseWriter) *traceResponseWriter {
	return &traceResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (w *traceResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func CreateTraceSimpleMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 从request获取reqeustid
			requestid := r.Header.Get("X-Request-ID")
			if requestid == "" {
				requestid = uuid.New().String()
			}
			atTime := time.Now()
			ctx := contextx.WithRequestContext(r.Context(), &contextx.RequestContext{
				TraceID: requestid,
				AtTime:  atTime,
			})
			wr := WrapResponseWriter(w)
			next.ServeHTTP(wr, r.WithContext(ctx))
			duration := time.Since(atTime)

			// 记录tarce日志
			traceLog, _ := json.Marshal(traceLogMessage{
				Level:      "",
				TraceID:    requestid,
				AtTime:     atTime.Format(time.DateTime),
				URL:        r.URL.String(),
				Method:     r.Method,
				Status:     wr.statusCode,
				DurationMs: duration.Milliseconds(),
			})
			logx.Core.Info("trace", string(traceLog))
		})
	}
}

type traceLogMessage struct {
	Level      string `json:"level"`
	TraceID    string `json:"trace_id"`
	AtTime     string `json:"at_time"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	Status     int    `json:"status"`
	DurationMs int64  `json:"duration_ms"`
}

// 实现logformatter, 并注册
type traceSimpleFormatter struct{}

func (f *traceSimpleFormatter) Format(level logx.LogLevel, message string) string {
	var data traceLogMessage

	// 如果message是一个json，还原，否则直接不变
	if json.Valid([]byte(message)) {
		if err := json.Unmarshal([]byte(message), &data); err == nil {
			data.Level = logx.GetLevelSTR(level)
			message = fmt.Sprintf("%+v", data)
		}
	}

	if json, err := json.Marshal(data); err == nil {
		return string(json)
	}

	return message
}

func init() {
	logx.RegisterFormatter("trace_simple", &traceSimpleFormatter{})
}
