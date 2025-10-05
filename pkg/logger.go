package pkg

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type CustomLogger struct {
	log       *slog.Logger
	tracer    trace.Tracer
	needTrace bool
}

type ErrorResponce struct {
	Logger  *CustomLogger
	Error   error
	Messege string
}

type InforResponce struct {
	Logger  *CustomLogger
	Messege string
}

func NewLogger(logger *slog.Logger, needTrace bool) *CustomLogger {
	return &CustomLogger{log: logger, needTrace: needTrace}
}

func (log *CustomLogger) SetTracer(tracer trace.Tracer) {
	log.tracer = tracer
}

func (log *CustomLogger) GetTracer() trace.Tracer {
	return log.tracer
}

func (log *CustomLogger) Info(info string, args ...any) *InforResponce {
	log.log.Info(info, args)
	return &InforResponce{Logger: log, Messege: info}
}

func (info *InforResponce) WithTrace(ctx context.Context) *InforResponce {
	if span := trace.SpanFromContext(ctx); info.Logger.needTrace && span.IsRecording() {
		span.SetStatus(codes.Ok, info.Messege)
	}
	return info
}

func (log *CustomLogger) Warn(msg string, err error) *ErrorResponce {
	log.log.Warn(msg, "error", err)
	return &ErrorResponce{Error: err, Messege: msg, Logger: log}
}

func (log *CustomLogger) Error(msg string, err error) *ErrorResponce {
	log.log.Error(msg, "error", err)
	return &ErrorResponce{Error: err, Messege: msg, Logger: log}
}

func (err *ErrorResponce) WithTrace(ctx context.Context) *ErrorResponce {
	if span := trace.SpanFromContext(ctx); err.Logger.needTrace && span.IsRecording() {
		span.RecordError(err.Error)
		span.SetStatus(codes.Error, err.Messege)
	}
	return err
}
