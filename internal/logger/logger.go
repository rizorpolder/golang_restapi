package logger

import (
	"fmt"
	"golang_restapi/internal/config"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
	"github.com/rs/zerolog"
)

func New(cfg *config.Base) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil || lvl == zerolog.NoLevel {
		lvl = zerolog.InfoLevel
	}

	var out io.Writer = os.Stdout // JSON по умолчанию
	if cfg.AppEnv == "local" || cfg.AppEnv == "development" {
		out = zerolog.ConsoleWriter{
			Out:         os.Stdout,
			TimeFormat:  "2006-01-02 15:04:05",
			FormatLevel: formatLevel,
		}
	}

	return zerolog.New(out).
		Level(lvl).
		With().
		Timestamp().
		Str("source", cfg.ServiceName).
		Str("env", cfg.AppEnv).
		Logger()
}

func formatLevel(i interface{}) string {
	if i == nil {
		return ""
	}
	switch level := strings.ToUpper(fmt.Sprintf("%s", i)); level {
	case "INFO":
		return color.New(color.FgGreen).Sprint("INF")
	case "WARN":
		return color.New(color.FgYellow).Sprint("WRN")
	case "ERROR":
		return color.New(color.FgRed).Sprint("ERR")
	case "DEBUG":
		return color.New(color.FgBlue).Sprint("DBG")
	default:
		return level
	}
}

var keyRe = regexp.MustCompile(`(?m)^(\s*)([A-Za-z0-9_]+):`)

func DumpConfig(log zerolog.Logger, cfg any) {
	spew.Config.Indent = " "
	spew.Config.DisablePointerAddresses = true
	spew.Config.DisableCapacities = true
	spew.Config.SortKeys = true

	keyColor := color.New(color.FgCyan).SprintfFunc()
	dump := keyRe.ReplaceAllStringFunc(spew.Sdump(cfg), func(s string) string {
		m := keyRe.FindStringSubmatch(s)
		if len(m) != 3 {
			return s
		}
		return fmt.Sprintf("%s%s:", m[1], keyColor(m[2]))
	})

	log.Debug().Msg("Loaded configuration\n" + dump)
}
