package log

import (
	"io"
	"log"
	"os"
	"sync"
)

var (
	errorLog = log.New(os.Stdout, "\033[31m[ERROR]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[32m[INFO]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{
		errorLog, infoLog,
	}
	mu sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

type LogLevel int

const (
	InfoLevel LogLevel = iota
	ErrorLevel
	Disabled
)

func SetLevel(level LogLevel) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	if level > ErrorLevel {
		errorLog.SetOutput(io.Discard)
	}
	if level > InfoLevel {
		infoLog.SetOutput(io.Discard)
	}
}
