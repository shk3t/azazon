package log

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var TestLogger *log.Logger
var TLog func(...any)

func Init(workDir string) {
	logDir := filepath.Join(workDir, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		panic("Can't create \"logs\" directory")
	}

	testLogFile, err := os.OpenFile(
		filepath.Join(logDir, "test.log"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		panic("Can't open \"test.log\" file")
	}

	TestLogger = log.New(testLogFile, "", log.LstdFlags|log.Lshortfile)
	TLog = TestLogger.Println
}

func Deinit() {
	writer := TestLogger.Writer()
	writeCloser, ok := writer.(io.WriteCloser)
	if ok {
		err := writeCloser.Close()
		if err != nil {
			panic("Can't close log file")
		}
	}
}