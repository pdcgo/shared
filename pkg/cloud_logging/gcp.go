package cloud_logging

import (
	"context"
	"log/slog"
	"os"
)

type CloudHandler struct {
	h slog.Handler
}

func (c *CloudHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return c.h.Enabled(ctx, lvl)
}

func (c *CloudHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs := []slog.Attr{
		slog.String("severity", levelToSeverity(r.Level)),
	}
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})
	newRec := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	newRec.AddAttrs(attrs...)
	return c.h.Handle(ctx, newRec)
}

func (c *CloudHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CloudHandler{h: c.h.WithAttrs(attrs)}
}

func (c *CloudHandler) WithGroup(name string) slog.Handler {
	return &CloudHandler{h: c.h.WithGroup(name)}
}

func levelToSeverity(l slog.Level) string {
	switch l {
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARNING"
	case slog.LevelError:
		return "ERROR"
	default:
		return "DEFAULT"
	}
}

func SetCloudLoggingDefault() {
	// Create base JSON handler
	base := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})

	// Wrap it with CloudHandler
	cloudHandler := &CloudHandler{h: base}

	// Set as default
	slog.SetDefault(slog.New(cloudHandler))
}
