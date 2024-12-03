package log

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

// levels
const (
	debugLevel   = 0
	releaseLevel = 1
	errorLevel   = 2
	fatalLevel   = 3
)

const (
	printDebugLevel   = "[debug] "
	printReleaseLevel = "[release] "
	printErrorLevel   = "[error] "
	printFatalLevel   = "[fatal] "
)

type Logger struct {
	level int

	// all
	baseLogger *log.Logger
	baseFile   *os.File

	// err
	errLogger *log.Logger
	errFile   *os.File
}

func New(strLevel string, pathname string, flag int) (*Logger, error) {
	// level
	var level int
	switch strings.ToLower(strLevel) {
	case "debug":
		level = debugLevel
	case "release":
		level = releaseLevel
	case "error":
		level = errorLevel
	case "fatal":
		level = fatalLevel
	default:
		return nil, errors.New("unknown level: " + strLevel)
	}

	// logger
	var baseLogger *log.Logger
	var baseFile *os.File

	var errBaseLogger *log.Logger
	var errBaseFile *os.File
	if pathname != "" {
		file, errFile, e := GetBaseFile(pathname)
		if e != nil {
			panic(e)
		}
		baseLogger = log.New(file, "", flag)
		baseFile = file

		errBaseLogger = log.New(errFile, "", flag)
		errBaseFile = errFile
	} else {
		baseLogger = log.New(os.Stdout, "", flag)
	}

	// new
	logger := new(Logger)
	logger.level = level
	logger.baseLogger = baseLogger
	logger.baseFile = baseFile
	logger.errLogger = errBaseLogger
	logger.errFile = errBaseFile
	return logger, nil
}

func RefreshLog(pathname string, logDays int) {
	go func() {
		f := func() {
			file, errFile, e := GetBaseFile(pathname)
			if e != nil {
				panic(e)
			}

			if logDays > 0 {
				deleteFilename := getFileName(time.Now().AddDate(0, 0, -logDays))
				deleteErrFilename := getErrFileName(time.Now().AddDate(0, 0, -logDays))
				os.Remove(path.Join(pathname, deleteFilename))
				os.Remove(path.Join(pathname, deleteErrFilename))
			}

			flag := log.LstdFlags | log.Llongfile
			gLogger.baseFile = file
			gLogger.baseLogger = log.New(file, "", flag)

			gLogger.errFile = errFile
			gLogger.errLogger = log.New(errFile, "", flag)
		}

		defer func() {
			if err := recover(); err != nil {
				gLogger.Error("%v,system will try again in 30 seconds", err)
				time.Sleep(30 * time.Second)
				f()
			}
		}()

		for {
			now := time.Now()
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now) + 5)
			<-t.C
			t.Stop()
			f()
		}
	}()
}

// It's dangerous to call the method on logging
func (logger *Logger) Close() {
	if logger.baseFile != nil {
		logger.baseFile.Close()
	}

	logger.baseLogger = nil
	logger.baseFile = nil
}

func (logger *Logger) doPrintf(level int, printLevel string, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	if logger.baseLogger == nil {
		panic("logger closed")
	}

	format = printLevel + format
	logger.baseLogger.Output(3, fmt.Sprintf(format, a...))

	if logger.errLogger != nil && level == errorLevel {
		logger.errLogger.Output(3, fmt.Sprintf(format, a...))
	}

	if level == fatalLevel {
		os.Exit(1)
	}
}

func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func (logger *Logger) Release(format string, a ...interface{}) {
	logger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func (logger *Logger) Error(format string, a ...interface{}) {
	logger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func (logger *Logger) Fatal(format string, a ...interface{}) {
	logger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

var gLogger, _ = New("debug", "", log.LstdFlags|log.Llongfile)

// It's dangerous to call the method on logging
func Export(logger *Logger) {
	if logger != nil {
		gLogger = logger
	}
}

func Debug(format string, a ...interface{}) {
	gLogger.doPrintf(debugLevel, printDebugLevel, format, a...)
}

func Release(format string, a ...interface{}) {
	gLogger.doPrintf(releaseLevel, printReleaseLevel, format, a...)
}

func Error(format string, a ...interface{}) {
	gLogger.doPrintf(errorLevel, printErrorLevel, format, a...)
}

func Fatal(format string, a ...interface{}) {
	gLogger.doPrintf(fatalLevel, printFatalLevel, format, a...)
}

func Close() {
	gLogger.Close()
}

func getFileName(t time.Time) string {
	return t.Format("2006-01-02") + ".log"
}

func getErrFileName(t time.Time) string {
	return t.Format("2006-01-02") + ".err.log"
}

func GetBaseFile(pathname string) (*os.File, *os.File, error) {

	if pathname != "" {
		// Create folder
		if _, err := os.Stat(pathname); err != nil {
			if os.IsNotExist(err) {
				if err := os.Mkdir(pathname, os.ModePerm); err != nil {
					panic(err)
				}
			}
			err = nil
		}

	} else {
		return nil, nil, fmt.Errorf("create log failed, log pathname is empty")
	}

	filename := path.Join(pathname, getFileName(time.Now()))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	errFilename := path.Join(pathname, getErrFileName(time.Now()))
	errFile, err := os.OpenFile(errFilename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, nil, err
	}

	return file, errFile, nil
}
