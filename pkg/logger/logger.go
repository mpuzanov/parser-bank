package logger

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogConf стуктура для настройки логирования
type LogConf struct {
	Level      string `yaml:"level" mapstructure:"level"`
	File       string `yaml:"file" mapstructure:"file"`
	FormatJSON bool   `yaml:"format_JSON" mapstructure:"format_JSON"`
}

// LogSugar logger по умолчанию
var LogSugar *zap.SugaredLogger

func init() {
	// иницилизируем logger по умолчанию
	LogSugar = zap.NewExample().Sugar()
}

func getZapLevel(level string) zap.AtomicLevel {
	switch level {
	case "info":
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "debug":
		return zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "error":
		return zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("02.01.2006 03:04:05 PM"))
}

//NewLogger Возвращаем инициализированный Logger
func NewLogger(config LogConf) *zap.Logger {
	EncodingFormat := "json"
	if !config.FormatJSON {
		EncodingFormat = "console"
	}

	OutputPath := []string{"stdout"}
	ErrorOutputPath := []string{"stderr"}

	if config.File != "" {
		_, err := os.Create(config.File)
		if err != nil {
			log.Printf("ошибка создания файла логов %s %v", config.File, err)
		} else {
			OutputPath = append(OutputPath, config.File)
			ErrorOutputPath = append(ErrorOutputPath, config.File)
		}
	}
	cfg := zap.Config{
		Encoding:         EncodingFormat,
		Level:            getZapLevel(config.Level),
		OutputPaths:      OutputPath,
		ErrorOutputPaths: ErrorOutputPath,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey: "time",
			//EncodeTime: zapcore.ISO8601TimeEncoder,
			EncodeTime: syslogTimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ := cfg.Build()
	return logger
}

//InitLogger Вариант инициализации логера
func InitLogger(config LogConf) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = syslogTimeEncoder //zapcore.ISO8601TimeEncoder
	cfg.CallerKey = "caller"
	cfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoder := zapcore.NewJSONEncoder(cfg)
	if !config.FormatJSON {
		encoder = zapcore.NewConsoleEncoder(cfg)
	}
	writerSyncer := zapcore.Lock(os.Stdout) // os.Stderr
	if config.File != "" {
		file, err := os.Create(config.File)
		if err != nil {
			log.Printf("ошибка создания файла логов %s %v", config.File, err)
		} else {
			writerSyncer = zapcore.NewMultiWriteSyncer(zapcore.Lock(file))
		}
	}
	level := getZapLevel(config.Level)
	logger := zap.New(zapcore.NewCore(encoder, writerSyncer, level))

	//zap.ReplaceGlobals(logger)
	return logger
}
