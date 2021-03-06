package common

import (
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log is the global logger.
var Log = zap.S()

// Version is the git commit hash.
var Version = ""

func init() {
	// set up a logger
	zcfg := zap.NewProductionConfig()
	zcfg.Encoding = "console"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zcfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zcfg.Level.SetLevel(zapcore.DebugLevel)

	logger, err := zcfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}

	Log = logger.Sugar()

	if Version == "" {
		Log.Info("Version not set, falling back to checking current directory.")

		git := exec.Command("git", "rev-parse", "--short", "HEAD")
		// ignoring errors *should* be fine? if there's no output we just fall back to "unknown"
		b, _ := git.Output()
		Version = strings.TrimSpace(string(b))
		if Version == "" {
			Version = "[unknown]"
		}
	}

	_, err = toml.DecodeFile("config.toml", &Conf)
	if err != nil {
		Log.Fatalf("Error reading configuration file: %v", err)
	}
}
