package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	timeFormat    = "[15:04:05.000]"
	minJSONLength = 2 // Minimum length for non-empty JSON object "{}"

	reset = "\033[0m"

	black        = 30
	red          = 31
	green        = 32
	yellow       = 33
	blue         = 34
	magenta      = 35
	cyan         = 36
	lightGray    = 37
	darkGray     = 90
	lightRed     = 91
	lightGreen   = 92
	lightYellow  = 93
	lightBlue    = 94
	lightMagenta = 95
	lightCyan    = 96
	white        = 97
)

func colorize(colorCode int, v string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), v, reset)
}

type Handler struct {
	h slog.Handler
	r func([]string, slog.Attr) slog.Attr
	b *bytes.Buffer
	m *sync.Mutex

	out io.Writer
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//nolint:exhaustruct // out field inherited from parent handler
	return &Handler{h: h.h.WithAttrs(attrs), b: h.b, r: h.r, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	//nolint:exhaustruct // out field inherited from parent handler
	return &Handler{h: h.h.WithGroup(name), b: h.b, r: h.r, m: h.m}
}

func (h *Handler) computeAttrs(ctx context.Context, r slog.Record) (map[string]any, error) {
	h.m.Lock()
	defer func() {
		h.b.Reset()
		h.m.Unlock()
	}()
	if err := h.h.Handle(ctx, r); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle: %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.b.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result: %w", err)
	}
	if _, ok := attrs["caller"]; !ok {
		var caller string
		_, file, line, ok := runtime.Caller(4)
		if ok {
			caller = file + ":" + strconv.Itoa(line)
		}
		attrs["caller"] = caller
	}
	return attrs, nil
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	var level string

	levelAttr := slog.Attr{
		Key:   slog.LevelKey,
		Value: slog.AnyValue(r.Level),
	}
	if h.r != nil {
		levelAttr = h.r([]string{}, levelAttr)
	}

	//nolint:exhaustruct // Empty Attr for comparison
	if !levelAttr.Equal(slog.Attr{}) {
		level = levelAttr.Value.String() + ":"

		if r.Level <= slog.LevelDebug {
			level = colorize(lightGray, level)
		} else if r.Level <= slog.LevelInfo {
			level = colorize(cyan, level)
		} else if r.Level < slog.LevelWarn {
			level = colorize(lightBlue, level)
		} else if r.Level < slog.LevelError {
			level = colorize(lightYellow, level)
		} else if r.Level <= slog.LevelError+1 {
			level = colorize(lightRed, level)
		} else if r.Level > slog.LevelError+1 {
			level = colorize(lightMagenta, level)
		}
	}

	var timestamp string

	timeAttr := slog.Attr{
		Key:   slog.TimeKey,
		Value: slog.StringValue(r.Time.Format(timeFormat)),
	}
	if h.r != nil {
		timeAttr = h.r([]string{}, timeAttr)
	}
	//nolint:exhaustruct // Empty Attr for comparison
	if !timeAttr.Equal(slog.Attr{}) {
		timestamp = colorize(lightGray, timeAttr.Value.String())
	}

	var msg string

	msgAttr := slog.Attr{
		Key:   slog.MessageKey,
		Value: slog.StringValue(r.Message),
	}
	if h.r != nil {
		msgAttr = h.r([]string{}, msgAttr)
	}
	//nolint:exhaustruct // Empty Attr for comparison
	if !msgAttr.Equal(slog.Attr{}) {
		msg = colorize(white, msgAttr.Value.String())
	}

	attrs, err := h.computeAttrs(ctx, r)
	if err != nil {
		return err
	}

	var caller string
	if c, ok := attrs["caller"]; ok {
		if cStr, ok := c.(string); ok {
			caller = colorize(yellow, cStr)
		}
		delete(attrs, "caller")
	}

	bytes, err := json.MarshalIndent(attrs, "", "  ")
	if err != nil {
		return fmt.Errorf("error when marshaling attrs: %w", err)
	}

	out := strings.Builder{}
	if len(timestamp) > 0 {
		out.WriteString(timestamp)
		out.WriteString(" ")
	}
	if len(level) > 0 {
		out.WriteString(level)
		out.WriteString(" ")
	}
	if len(msg) > 0 {
		out.WriteString(msg)
		out.WriteString(" ")
	}
	if len(caller) > 0 {
		out.WriteString(colorize(yellow, "CALLER: "))
		out.WriteString(caller)
		out.WriteString(" ")
	}
	if len(bytes) > minJSONLength {
		out.WriteString(colorize(darkGray, string(bytes)))
	}
	out.WriteString("\n")

	if _, err = h.out.Write([]byte(out.String())); err != nil {
		return fmt.Errorf("error writing to output: %w", err)
	}

	return nil
}

func suppressDefaults(
	next func([]string, slog.Attr) slog.Attr,
) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			//nolint:exhaustruct // Intentionally returning empty Attr to suppress default fields
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}

func NewHandler(out io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		//nolint:exhaustruct // Empty options uses slog defaults
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}

	return &Handler{
		b: b,

		h: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		r:   opts.ReplaceAttr,
		m:   &sync.Mutex{},
		out: out,
	}
}
