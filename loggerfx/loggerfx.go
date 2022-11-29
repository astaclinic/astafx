package loggerfx

import (
	"fmt"
	"os"
	"path"

	"github.com/fatih/color"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Options(
	fx.Provide(New),
	fx.WithLogger(func(logger *zap.SugaredLogger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: logger.Desugar()}
	}),
)

type LoggerConfig struct {
	Path string `mapstructure:"path" validate:"required,dir"`
}

func New(config *LoggerConfig) (*zap.SugaredLogger, error) {
	// create directory if needed
	err := os.MkdirAll(config.Path, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("error in creating log file folder for writing: %w", err)
	}

	// open a new file
	logFile, err := os.OpenFile(path.Join(config.Path, "latest.log"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("error in opening log file for writing: %w", err)
	}

	// setting the log level for file/console log output
	fileLogLevel := zapcore.InfoLevel // enable if level >= info, can change to any predicate accepting a level argument
	consoleLogLevel := zapcore.InfoLevel

	// setup the encoders
	fileEncoderConfig := zap.NewProductionEncoderConfig()
	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoderConfig := zap.NewProductionEncoderConfig()
	colorMap := map[zapcore.Level]*color.Color{
		zapcore.DebugLevel:  color.New(color.FgMagenta),
		zapcore.InfoLevel:   color.New(color.FgBlue),
		zapcore.WarnLevel:   color.New(color.FgYellow),
		zapcore.ErrorLevel:  color.New(color.FgRed),
		zapcore.DPanicLevel: color.New(color.FgRed, color.Bold),
		zapcore.FatalLevel:  color.New(color.FgRed, color.Bold),
		zapcore.PanicLevel:  color.New(color.FgRed, color.Bold),
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
		zapcore.NewCore(fileEncoder, zapcore.Lock(logFile), fileLogLevel),
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), consoleLogLevel),
	)

	return zap.New(core, zap.AddCaller()).Sugar(), nil
}
