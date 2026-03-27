package log

import (
	"os"
	"path/filepath"

	"dario.cat/mergo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	logger *zap.Logger
	config Config
	hooks  []Hook
}

type Config struct {
	Name        string        `json:"name"`
	Debug       bool          `json:"debug"`
	SkipLevel   int           `json:"skip_level"`
	Level       zapcore.Level `json:"level"`
	LevelKey    string        `json:"level_key"`
	TimeKey     string        `json:"time_key"`
	CallerKey   string        `json:"caller_key"`
	FunctionKey string        `json:"function_key"`
	NameKey     string        `json:"name_key"`
	Encoding    string        `json:"encoding"` // zap encoder
	Includes    []string      `json:"includes"`
	Excludes    []string      `json:"excludes"`

	// Output controls where logs are written: "stdio" or "file" (default: file)
	Output string `json:"output"`
	// File holds file-based logging configuration
	File FileConfig `json:"file"`
}

// FileConfig holds file-based logging options (lumberjack v2).
type FileConfig struct {
	Path       string `json:"path"`
	MaxSize    int    `json:"max_size"`    // megabytes
	MaxAge     int    `json:"max_age"`     // days
	MaxBackups int    `json:"max_backups"` // files
	LocalTime  bool   `json:"local_time"`
}

type (
	ObjectEncoder   = zapcore.ObjectEncoder
	ObjectMarshaler = zapcore.ObjectMarshaler
)

var (
	// String     = zap.String
	// Bool       = zap.Bool
	// Strings    = zap.Strings
	// ByteString = zap.ByteString
	// Float64    = zap.Float64
	// Int64      = zap.Int64
	// Int32      = zap.Int32
	// Int        = zap.Int
	// Uint       = zap.Uint
	// Uint64     = zap.Uint64
	// Duration   = zap.Duration
	// Object     = zap.Object
	// Namespace  = zap.Namespace
	// Reflect    = zap.Reflect
	// Stack      = zap.Stack
	// Time       = zap.Time
	// Skip       = zap.Skip()

	Cause = func(err error) zap.Field {
		return NamedError("error", err)
	}

	NamedError = func(key string, err error) zap.Field {
		if err == nil {
			return zap.Skip()
		} else {
			return Any(key, err)
		}
	}

	Any = func(key string, value any) zap.Field {
		if value == nil {
			return zap.Skip()
		}
		return zap.Any(key, value)
	}

	EncodeStringSlice = func(lines []string) zapcore.ArrayMarshalerFunc {
		return func(encoder zapcore.ArrayEncoder) error {
			for _, line := range lines {
				encoder.AppendString(line)
			}
			return nil
		}
	}
)

var defaultConfig = Config{
	Name:      "default",
	Debug:     false,
	SkipLevel: 1,
	Level:     zap.InfoLevel,
	LevelKey:  "level",
	TimeKey:   "time",
	CallerKey: "label",
	NameKey:   "logger",
	Encoding:  "json",
	Includes:  []string{},
	Excludes:  []string{},
	Output:    "stdio",
	File: FileConfig{
		Path:       "logs/x2o.log",
		MaxSize:    50,
		MaxAge:     30,
		MaxBackups: 10,
		LocalTime:  true,
	},
}

func New(config Config) *Logger {
	err := mergo.Merge(&config, defaultConfig)
	if err != nil {
		panic(err)
	}

	encCfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       config.LevelKey,
		TimeKey:        config.TimeKey,
		NameKey:        config.NameKey,
		CallerKey:      config.CallerKey,
		FunctionKey:    config.FunctionKey,
		StacktraceKey:  "stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Select encoder based on config.Encoding
	var encoder zapcore.Encoder

	switch config.Encoding {
	case "json", "":
		encoder = zapcore.NewJSONEncoder(encCfg)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encCfg)
	case "console_json":
		encoder = NewConsoleJSONEncoder(encCfg)
	default:
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	// Select writer syncer based on output target
	var ws zapcore.WriteSyncer

	switch config.Output {
	case "file":
		path := config.File.Path
		if path == "" {
			path = "logs/axonhub.log"
		}

		if dir := filepath.Dir(path); dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o750); err != nil {
				panic(err)
			}
		}

		lj := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    config.File.MaxSize,
			MaxAge:     config.File.MaxAge,
			MaxBackups: config.File.MaxBackups,
			LocalTime:  config.File.LocalTime,
		}
		ws = zapcore.AddSync(lj)
	case "stdio", "stdout", "console", "":
		ws = zapcore.AddSync(os.Stdout)
	default:
		ws = zapcore.AddSync(os.Stdout)
	}

	// Build core with DebugLevel enabler; per-level gating is handled by Logger methods.
	core := zapcore.NewCore(encoder, ws, zapcore.DebugLevel)

	opts := []zap.Option{zap.AddStacktrace(zapcore.DPanicLevel), zap.ErrorOutput(zapcore.AddSync(os.Stderr))}
	if config.SkipLevel != 0 {
		opts = append(opts, zap.AddCallerSkip(config.SkipLevel))
	} else {
		opts = append(opts, zap.AddCallerSkip(defaultConfig.SkipLevel))
	}

	if len(config.Includes) > 0 || len(config.Excludes) > 0 {
		opts = append(opts, withNameFilter(config.Includes, config.Excludes))
	}

	zapLogger := zap.New(core, opts...).Named(config.Name)

	return &Logger{
		config: config,
		logger: zapLogger,
		hooks:  []Hook{HookFunc(contextFields)},
	}
}
