package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	NONE = iota
	ERROR
	DEBUG
)

type Logger struct {
	erorLogger  *log.Logger
	debugLogger *log.Logger
	printLogs   bool
	logLevel    int
}

func (l *Logger) Debug(log_entry string) {
	if l.logLevel >= DEBUG {
		l.printLogsToConsole(log_entry)
		l.debugLogger.Print(log_entry)
	}
}

func (l *Logger) Debugf(format_string string, args ...any) {
	if l.logLevel >= DEBUG {
		l.printfLogsToConsole(format_string, args...)
		l.debugLogger.Printf(format_string, args...)
	}
}

func (l *Logger) Error(log_entry string) {
	if l.logLevel >= ERROR {
		l.printLogsToConsole(log_entry)
		l.erorLogger.Print(log_entry)
	}
}

func (l *Logger) Errorf(format_string string, args ...any) {
	if l.logLevel >= ERROR {
		l.printfLogsToConsole(format_string, args...)
		l.erorLogger.Printf(format_string, args...)
	}
}

func (l *Logger) printLogsToConsole(log_entry string) {
	if l.printLogs {
		log.Print(log_entry)
	}
}

func (l *Logger) printfLogsToConsole(format string, rest ...any) {
	if l.printLogs {
		log.Printf(format, rest...)
	}
}

func (l *Logger) setLogLevel(level int) {
	if (level < NONE) || (level > DEBUG) {
		l.Error(fmt.Sprintf("Cannot set level to %d", level))
		return
	}

	l.logLevel = level
}

func openLogFile(path string, flags int, mode int) *os.File {
	logfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Cannot open error %s.", path)
	}
	return logfile
}

func newLogger(printLogs bool) *Logger {
	logPath := os.TempDir()
	flags := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	mode := 0600
	error_logfile := openLogFile(filepath.Join(logPath, "error.log"), flags, mode)
	debug_logfile := openLogFile(filepath.Join(logPath, "debug.log"), flags, mode)

	return &Logger{
		erorLogger:  log.New(error_logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(debug_logfile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		printLogs:   printLogs,
	}
}

var LOGGER *Logger = newLogger(false)

func Debug(log_entry string) {
	LOGGER.Debug(log_entry)
}

func Debugf(format_string string, args ...any) {
	LOGGER.Debugf(format_string, args...)
}

func Error(log_entry string) {
	LOGGER.Error(log_entry)
}

func Errorf(format_string string, args ...any) {
	LOGGER.Errorf(format_string, args...)
}

func Fatal(log_entry string) {
	log.Fatal(log_entry)
}

func Fatalf(format_string string, args ...any) {
	log.Fatalf(format_string, args...)
}

func SetLogLevel(log_level int) {
	LOGGER.setLogLevel(log_level)
}

func SetLogPrint(print bool) {
	LOGGER.printLogs = print
}
