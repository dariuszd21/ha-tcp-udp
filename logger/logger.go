package logger

import (
	"path/filepath"
	"log"
	"os"
)

type Logger struct {
	erorLogger *log.Logger
	debugLogger *log.Logger
	printLogs bool
}

func (l* Logger) Debug(log_entry string) {
	if l.printLogs {
		log.Print(log_entry)
	}
	l.debugLogger.Print(log_entry)
}

func (l* Logger) Error(log_entry string) {
	if l.printLogs {
		log.Print(log_entry)
	}
	l.erorLogger.Print(log_entry)
}

func openLogFile(path string, flags int, mode int) *os.File {
	logfile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Cannot open error %s.", path)
	}
	return logfile
}

func NewLogger(printLogs bool) *Logger {
	// TODO: Read from configuration file
	logPath := os.TempDir()
	flags := os.O_APPEND|os.O_CREATE|os.O_WRONLY 
	mode := 0600
	error_logfile := openLogFile(filepath.Join(logPath, "error.log"), flags, mode)
	debug_logfile := openLogFile(filepath.Join(logPath, "debug.log"), flags, mode)

	return &Logger{
		erorLogger: log.New(error_logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		debugLogger: log.New(debug_logfile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}


