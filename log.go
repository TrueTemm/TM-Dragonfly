/*
 _____               _____
|_   _| __ _   _  __|_   _|__ _ __ ___  _ __ ___
  | || '__| | | |/ _ \| |/ _ \ '_ ` _ \| '_ ` _ \
  | || |  | |_| |  __/| |  __/ | | | | | | | | | |
  |_||_|   \__,_|\___||_|\___|_| |_| |_|_| |_| |_|

 _____ __  __       ____                               __ _
|_   _|  \/  |     |  _ \ _ __ __ _  __ _  ___  _ __  / _| |_   _
  | | | |\/| |_____| | | | '__/ _` |/ _` |/ _ \| '_ \| |_| | | | |
  | | | |  | |_____| |_| | | | (_| | (_| | (_) | | | |  _| | |_| |
  |_| |_|  |_|     |____/|_|  \__,_|\__, |\___/|_| |_|_| |_|\__, |
                                    |___/                   |___/

@author TrueTemm
@link   https://github.com/TrueTemm
TM-Dragonfly Project
*/

package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"
)

const (
	cReset   = "\x1b[0m"
	cGray    = "\x1b[90m"
	cRed     = "\x1b[31m"
	cGreen   = "\x1b[32m"
	cYellow  = "\x1b[33m"
	cBlue    = "\x1b[34m"
	cMagenta = "\x1b[35m"
	cCyan    = "\x1b[36m"
	cBold    = "\x1b[1m"

	cGold = "\x1b[38;5;220m"
	cLime = "\x1b[38;5;120m" // tm-dragonfly brand colour — light green
)

type tmHandler struct {
	w     io.Writer
	mu    *sync.Mutex
	level slog.Level
	attrs []slog.Attr
}

func newTMHandler(w io.Writer, level slog.Level) *tmHandler {
	return &tmHandler{w: w, mu: &sync.Mutex{}, level: level}
}

func (h *tmHandler) Enabled(_ context.Context, l slog.Level) bool { return l >= h.level }

func (h *tmHandler) levelColour(l slog.Level) (string, string) {
	switch {
	case l >= slog.LevelError:
		return "ERROR", cRed
	case l >= slog.LevelWarn:
		return "WARN", cYellow
	case l >= slog.LevelInfo:
		return "INFO", cLime
	default:
		return "DEBUG", cGray
	}
}

func (h *tmHandler) Handle(_ context.Context, r slog.Record) error {
	name, colour := h.levelColour(r.Level)

	var b strings.Builder
	b.WriteString(cGray + r.Time.Format("15:04:05") + cReset + " ")
	b.WriteString(cLime + cBold + "TM-Dragonfly" + cReset + " ")
	b.WriteString(colour + "[" + name + "]" + cReset + " ")
	b.WriteString(r.Message)

	writeAttr := func(a slog.Attr) {
		b.WriteString(fmt.Sprintf("  %s%s%s=%v", cCyan, a.Key, cReset, a.Value.Any()))
	}
	for _, a := range h.attrs {
		writeAttr(a)
	}
	r.Attrs(func(a slog.Attr) bool {
		writeAttr(a)
		return true
	})
	b.WriteString("\n")

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := io.WriteString(h.w, b.String())
	return err
}

func (h *tmHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	n := *h
	n.attrs = append(append([]slog.Attr{}, h.attrs...), attrs...)
	return &n
}

func (h *tmHandler) WithGroup(string) slog.Handler { return h }

var _ slog.Handler = (*tmHandler)(nil)

func colourLogger(w io.Writer, level slog.Level) *slog.Logger {
	return slog.New(newTMHandler(w, level))
}
