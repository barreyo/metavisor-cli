package logging

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type level int

const (
	// LevelDebug logs everything
	LevelDebug level = iota
	// LevelInfo only logs Info and above
	LevelInfo
	// LevelWarning only logs Warning and above
	LevelWarning
	// LevelError only logs Error and Output
	LevelError
	// LevelOutput only logs output messages, so no logs really
	LevelOutput

	templateDebug   = "DEBUG: %v"
	templateInfo    = "INFO:  %v"
	templateWarning = "WARN:  %v"
	templateError   = "ERROR: %v"
	templateOutput  = "OUTPUT: %s"
	templateFatal   = "FATAL: %s"
)

var (
	// LogLevel sets the minimum log level that will be shown in console
	LogLevel = LevelInfo
	// LogFileNamePrefix is the prefix of the log file written
	LogFileNamePrefix = "log"
	// LogToFile determines whether to also log output to a file or not
	LogToFile = true

	outLogger   = log.New(os.Stdout, "", 0)
	termLogger  = log.New(os.Stderr, "", 0)
	fileLog     *log.Logger
	logFilePath = ""
)

func LogFile() string {
	return logFilePath
}

func Debug(v ...interface{}) {
	print(LevelDebug, templateDebug, v...)
}

func Debugf(t string, v ...interface{}) {
	print(LevelDebug, templateDebug, fmt.Sprintf(t, v...))
}

func Info(v ...interface{}) {
	print(LevelInfo, templateInfo, v...)
}

func Infof(t string, v ...interface{}) {
	print(LevelInfo, templateInfo, fmt.Sprintf(t, v...))
}

func Warning(v ...interface{}) {
	print(LevelWarning, templateWarning, v...)
}

func Warningf(t string, v ...interface{}) {
	print(LevelWarning, templateWarning, fmt.Sprintf(t, v...))
}

func Error(v ...interface{}) {
	print(LevelError, templateError, v...)
}

func Errorf(t string, v ...interface{}) {
	print(LevelError, templateError, fmt.Sprintf(t, v...))
}

func Output(v ...interface{}) {
	outLogger.Println(v...)
	if fileLogger() != nil {
		fileLogger().Printf(templateOutput, fmt.Sprintln(v...))
	}
}

func Outputf(t string, v ...interface{}) {
	outLogger.Printf(t, v...)
	if fileLogger() != nil {
		fileLogger().Printf(templateOutput, fmt.Sprintf(t, v...))
	}
}

func Fatal(v ...interface{}) {
	termLogger.Printf(templateFatal, fmt.Sprintln(v...))
	if fileLogger() != nil {
		fileLogger().Printf(templateFatal, fmt.Sprintln(v...))
		outLogger.Printf("Logs are available at:\n%s", logFilePath)
	}
	os.Exit(1)
}

func Fatalf(t string, v ...interface{}) {
	termLogger.Printf(templateFatal, fmt.Sprintf(t, v...))
	if fileLogger() != nil {
		fileLogger().Printf(templateFatal, fmt.Sprintf(t, v...))
		outLogger.Printf("Logs are available at:\n%s", logFilePath)
	}
	os.Exit(1)
}

func print(lvl level, template string, v ...interface{}) {
	if lvl >= LogLevel {
		if lvl == LevelInfo && LogLevel != LevelDebug {
			termLogger.Print(fmt.Sprintln(v...))
		} else {
			termLogger.Printf(template, fmt.Sprintln(v...))
		}
	}
	if fileLogger() != nil {
		fileLogger().Printf(template, fmt.Sprintln(v...))
	}
}

func fileLogger() *log.Logger {
	if !LogToFile {
		return nil
	}
	if fileLog != nil {
		return fileLog
	}

	f, err := ioutil.TempFile("", LogFileNamePrefix)
	if err != nil {
		Errorf("Could not write log to file: %s", err)
		return nil
	}
	fileLog = log.New(f, "", log.Ldate|log.Ltime)
	logFilePath = f.Name()
	return fileLog
}
