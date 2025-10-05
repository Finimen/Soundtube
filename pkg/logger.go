package pkg

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CustomLogger struct {
	log    *slog.Logger
	tracer trace.Tracer
}

func NewLogger(logger *slog.Logger) CustomLogger {
	return CustomLogger{log: logger}
}

func (log *CustomLogger) SetTracer(tracer trace.Tracer) {
	log.tracer = tracer
}

func (log *CustomLogger) Info(info string, args ...any) {
	log.log.Info(info, args)
}

func (log *CustomLogger) Error(ctx context.Context, msg string, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, msg)
	log.log.Error(msg, "error", err, "trace_id", span.SpanContext().TraceID().String())
}
