package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type LogLevel string

var (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
	FatalLevel LogLevel = "FATAL"
)

type LogColor *color.Color

var (
	DebugColor LogColor = color.New(color.FgMagenta)
	InfoColor  LogColor = color.New(color.FgBlue)
	WarnColor  LogColor = color.New(color.FgYellow)
	ErrorColor LogColor = color.New(color.FgRed)
	FatalColor LogColor = color.New(color.FgRed, color.Bold)
)

var LogLevelColor = map[LogLevel]LogColor{
	DebugLevel: DebugColor,
	InfoLevel:  InfoColor,
	WarnLevel:  WarnColor,
	ErrorLevel: ErrorColor,
	FatalLevel: FatalColor,
}

func Log(logLevel LogLevel, message string) error {
	_, err := fmt.Printf("%s\t%s\t%s\n",
		time.Now().Format(time.RFC3339),
		(*color.Color)(LogLevelColor[logLevel]).Sprintf("[%s]", logLevel),
		message,
	)
	return err
}

func Info(message string) error {
	return Log(InfoLevel, message)
}

func Infof(format string, a ...any) error {
	return Info(fmt.Sprintf(format, a...))
}

func Warn(message string) error {
	return Log(WarnLevel, message)
}

func Warnf(format string, a ...any) error {
	return Warn(fmt.Sprintf(format, a...))
}

func Fatal(message string) {
	Log(FatalLevel, message)
	os.Exit(1)
}

func Fatalf(format string, a ...any) {
	Fatal(fmt.Sprintf(format, a...))
	os.Exit(1)
}
