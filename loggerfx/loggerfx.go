package loggerfx

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/astaclinic/astafx/config"
	"github.com/astaclinic/astafx/logger"
)

var Module = fx.Options(
	fx.Provide(New),
	fx.WithLogger(func(logger *zap.SugaredLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Desugar()}
	}),
)

type LoggerConfig struct {
	Path string `mapstructure:"path" yaml:"path" validate:"required"`
}

func init() {
	// config must have a default value for viper to load config from env variables
	// default value of empty string (zero value) will not pass the "required" config validation
	viper.SetDefault("logs.path", path.Join("/var/log", config.GetPackageName()))
}

func New(config *LoggerConfig) (*zap.SugaredLogger, error) {
	// create directory if needed
	err := os.MkdirAll(config.Path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error in creating log file folder for writing: %w", err)
	}

	// create a new writer for log rotation
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename: path.Join(config.Path, "server.log"),
	})

	// setting the log level for file/console log output
	fileLogLevel := zapcore.InfoLevel // enable if level >= info, can change to any predicate accepting a level argument
	consoleLogLevel := zapcore.InfoLevel

	// setup the encoders
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	colorMap := map[zapcore.Level]*color.Color{
		zapcore.DebugLevel:  logger.DebugColor,
		zapcore.InfoLevel:   logger.InfoColor,
		zapcore.WarnLevel:   logger.WarnColor,
		zapcore.ErrorLevel:  logger.ErrorColor,
		zapcore.DPanicLevel: logger.FatalColor,
		zapcore.FatalLevel:  logger.FatalColor,
		zapcore.PanicLevel:  logger.FatalColor,
	}
	consoleEncoderConfig.EncodeLevel = func(l zapcore.Level, pae zapcore.PrimitiveArrayEncoder) {
		// custom encoding of level string as [INFO] style
		pae.AppendString(colorMap[l].Sprintf("[%s]", l.CapitalString()))
	}
	consoleEncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	consoleEncoderConfig.EncodeCaller = func(ec zapcore.EntryCaller, pae zapcore.PrimitiveArrayEncoder) {
		// custom encoding of the caller, now is set to the trimmed file path
		pae.AppendString(ec.TrimmedPath())
	}
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// create the two cores for the logger
	// when writing to a file, the *os.File need to be locked with Lock() for concurrent access
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileWriter, fileLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLogLevel),
	)

	return zap.New(core, zap.AddCaller()).Sugar(), nil
}
