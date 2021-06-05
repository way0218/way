package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	DefaultLevel      = zapcore.InfoLevel
	DefaultTimeLayout = "2006-01-02 15:04:05"
	DefaultMaxSize    = 1024
	DefaultMaxBackups = 100
	DefaultMaxAge     = 100
)

// Option setup config
type Option func(opt *option)

type option struct {
	level          zapcore.Level //日志等级
	file           string
	timeLayout     string // 日志时间输出格式
	disableConsole bool
	maxSize        int  // 日志文件大小(MB)
	maxBackups     int  // 最多存在多少个文件
	maxAge         int  // 最多保存多少天
	compress       bool // 日志是否压缩
	consoleEncoder bool // console 日志格式
}

func WithDebugLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.DebugLevel
	}
}

func WithInfoLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.InfoLevel
	}
}

func WithWarnLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.WarnLevel
	}
}

func WithErrorLevel() Option {
	return func(opt *option) {
		opt.level = zapcore.ErrorLevel
	}
}

func WithTimeLayout(timeLayout string) Option {
	return func(opt *option) {
		opt.timeLayout = timeLayout
	}
}

func WithDisableConsole() Option {
	return func(opt *option) {
		opt.disableConsole = true
	}
}

func WithCompress() Option {
	return func(opt *option) {
		opt.compress = true
	}
}

func WithFileRotation(file string) Option {
	return func(opt *option) {
		opt.file = file
	}
}

func WithConsoleEncoder() Option {
	return func(opt *option) {
		opt.consoleEncoder = true
	}
}

func WithMaxSize(maxSize int) Option {
	return func(opt *option) {
		opt.maxSize = maxSize
	}
}

func WithMaxAge(maxAge int) Option {
	return func(opt *option) {
		opt.maxAge = maxAge
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(opt *option) {
		opt.maxBackups = maxBackups
	}
}

func NewLogger(opts ...Option) *zap.Logger {
	opt := &option{
		level:          DefaultLevel,
		timeLayout:     DefaultTimeLayout,
		disableConsole: false,
		maxSize:        DefaultMaxSize,
		maxBackups:     DefaultMaxBackups,
		maxAge:         DefaultMaxAge,
		compress:       false,
	}
	for _, f := range opts {
		f(opt)
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "line",
		MessageKey:       "msg",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.LowercaseLevelEncoder,               // 小写编码器
		EncodeTime:       zapcore.TimeEncoderOfLayout(opt.timeLayout), // 时间格式
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	if opt.consoleEncoder {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()
	if !opt.disableConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder,
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(encoder,
				zapcore.NewMultiWriteSyncer(stderr),
				highPriority,
			),
		)
	}
	if opt.file != "" {
		core = zapcore.NewTee(core, zapcore.NewCore(encoder,
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   opt.file,
				MaxSize:    opt.maxSize,
				MaxAge:     opt.maxAge,
				MaxBackups: opt.maxBackups,
				Compress:   opt.compress,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= opt.level
			}),
		),
		)
	}
	logger := zap.New(core, zap.AddCaller(), zap.ErrorOutput(stderr))
	return logger
}
