package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var Loggers = struct {
	Debug *log.Logger
	Event *log.Logger
	Test  *log.Logger
	Run   *log.Logger
}{}

var fileLoggers = map[string]**log.Logger{
	"debug.log": &Loggers.Debug,
	"event.log": &Loggers.Event,
}

var (
	Debug func(...any)
)

func Init(workDir string) error {
	logDir := filepath.Join(workDir, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	for fileName, loggerPtr := range fileLoggers {
		logFile, err := os.OpenFile(
			filepath.Join(logDir, fileName),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0644,
		)
		if err != nil {
			return err
		}

		*loggerPtr = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
	}

	Loggers.Run = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	Loggers.Test = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)
	Debug = Loggers.Debug.Println

	return nil
}

func Deinit() {
	for _, loggerPtr := range fileLoggers {
		if loggerPtr == nil || *loggerPtr == nil {
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