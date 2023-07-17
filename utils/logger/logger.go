package logger

import (
	"fmt"
	"moco/utils/console/color"
	"time"
)

// Logger struct
type Logger struct {
	// Context string
	Context string
}

// NewLogger returns a new Logger
func New(context string) *Logger {
	return &Logger{
		Context: context,
	}
}

// Info logs a message with INFO level
func (l *Logger) Info(message ...string) {
	l.log("INFO", color.Green, message...)
}

// Infof logs a message with INFO level
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log("INFO", color.Green, fmt.Sprintf(format, args...))
}

// Error logs a message with ERROR level
func (l *Logger) Error(message ...string) {
	l.log("ERROR", color.Red, message...)
}

func (l *Logger) Errore(err error) {
	l.log("ERROR", color.Red, err.Error())
}

// Errorf logs a message with ERROR level
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log("ERROR", color.Red, fmt.Sprintf(format, args...))
}

// Warn logs a message with WARN level
func (l *Logger) Warn(message ...string) {
	l.log("WARN", color.Yellow, message...)
}

// Warnf logs a message with WARN level
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log("WARN", color.Yellow, fmt.Sprintf(format, args...))
}

// Debug logs a message with DEBUG level
func (l *Logger) Debug(message ...string) {
	l.log("DEBUG", color.Purple, message...)
}

// Debugf logs a message with DEBUG level
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log("DEBUG", color.Purple, fmt.Sprintf(format, args...))
}

// Panic logs a message with PANIC level
func (l *Logger) Panic(message ...string) {
	l.log("PANIC", color.Red, message...)
}

// Panice logs a message with PANIC level
func (l *Logger) Panice(err error) {
	l.log("PANIC", color.Red, err.Error())
}

// Panicf logs a message with PANIC level
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.log("PANIC", color.Red, fmt.Sprintf(format, args...))
}

func (l *Logger) log(level string, c string, message ...string) {
	// join all messages separated by space
	messageStr := ""
	for _, m := range message {
		messageStr += m + " "
	}

	fmt.Printf(
		"%s %s%s%s [%s]%s %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		c,
		fmt.Sprintf("%-5s", level),
		color.Blue,
		l.Context,
		color.Reset,
		messageStr,
	)
}
