package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

var (
	LoggingActive = false
	LogService    *AsyncLogger
)

func StartLogging() error {
	logS, err := New(time.Now().Format("2006-01-02") + "_events.log")
	if err != nil {
		return err
	}
	LoggingActive = true
	LogService = logS
	return nil
}

type AsyncLogger struct {
	file *os.File
	ch   chan string
	wg   *sync.WaitGroup
}

func New(filename string) (*AsyncLogger, error) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("error al abrir el archivo de log: %v", err)
	}
	logger := &AsyncLogger{
		file: file,
		ch:   make(chan string, 100),
		wg:   new(sync.WaitGroup),
	}

	logger.wg.Add(1)
	go logger.runLogger()

	return logger, nil
}

func (l *AsyncLogger) runLogger() {
	defer l.wg.Done()

	// Hasta que el canal l .ch sea cerrado
	for message := range l.ch {
		fmt.Fprintln(l.file, message)
	}
}

func SendLog(format string, v ...any) {
	if !LoggingActive {
		return
	}
	message := fmt.Sprintf(format, v...)
	LogService.Log(message)
}

func (l *AsyncLogger) Log(message string) {

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	logTime := time.Now().Format("2006/01/02 15:04:05")
	logMsg := fmt.Sprintf("%s %s:%d: %s", logTime, file, line, message)

	select {
	case l.ch <- logMsg:
	default:
		fmt.Println("ADVERTENCIA: Buffer de logs lleno. Mensaje perdido.")
	}
}

func (l *AsyncLogger) Logf(format string, v ...any) {
	message := fmt.Sprintf(format, v...)

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}

	logTime := time.Now().Format("2006/01/02 15:04:05")
	logMsg := fmt.Sprintf("%s %s:%d: %s", logTime, file, line, message)

	select {
	case l.ch <- logMsg:
	default:
		fmt.Println("ADVERTENCIA: Buffer de logs lleno. Mensaje perdido.")
	}
}

func (l *AsyncLogger) Close() {
	close(l.ch)
	l.wg.Wait()
	l.file.Close()
}
