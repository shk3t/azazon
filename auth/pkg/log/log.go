package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var Loggers = struct {
	Request *log.Logger
	Test    *log.Logger
}{}

var allLoggers = map[string]**log.Logger{
	"request.log": &Loggers.Request,
	"test.log":    &Loggers.Test,
}

var TLog func(...any)

func Init(workDir string) {
	logDir := filepath.Join(workDir, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		panic("Can't create \"logs\" directory")
	}

	for fileName, loggerPtr := range allLoggers {
		logFile, err := os.OpenFile(
			filepath.Join(logDir, fileName),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			panic("Can't open \"test.log\" file")
		}

		*loggerPtr = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	}

	TLog = Loggers.Test.Println
}

func Deinit() {
	for _, loggerPtr := range allLoggers {
		if loggerPtr == nil {
			continue
		}

		writer := (*loggerPtr).Writer()
		writeCloser, ok := writer.(io.WriteCloser)
		if ok {
			err := writeCloser.Close()
			if err != nil {
				panic("Can't close log file")
			}
		}
	}
}