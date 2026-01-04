package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var (
	appLog   *log.Logger
	appFile  *os.File
	logMutex sync.Mutex
)

func init() {
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	
	logDir := filepath.Dir(execPath)
	logPath := filepath.Join(logDir, "app.log")
	appFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		appFile = os.Stderr
	}
	
	appLog = log.New(appFile, "", log.LstdFlags)
}

func Log(message string) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Println(message)
}

func Logf(format string, args ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	appLog.Printf(format, args...)
}

func CloseLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()
	if appFile != nil && appFile != os.Stderr {
		appFile.Close()
	}
}

func Print(message string) {
	fmt.Print(message)
}

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
