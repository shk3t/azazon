package service

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ReadEmails(recipient string) ([]string, error) {
	stat, err := mailer.logFile.Stat()
	if err != nil {
		return nil, err
	}

	buf := make([]byte, stat.Size())
	_, err = mailer.logFile.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	allMessages := strings.Split(string(buf), "\n")
	filteredMessages := []string{}
	for _, msg := range allMessages {
		if strings.Contains(msg, recipient) {
			filteredMessages = append(filteredMessages, msg)
		}
	}

	return filteredMessages, nil
}

func SendEmail(recipient string, message string) error {
	if mailer.logger == nil {
		return errors.New("Mailer is not inited")
	}

	mailer.logger.Printf("Message for %s: %s", recipient, message)
	return nil
}

func FmtUserById(userId int) string {
	return fmt.Sprintf("<User_%d>", userId)
}

var mailer = struct {
	logger  *log.Logger
	logFile *os.File
}{}

func InitMailer(workDir string) error {
	logDir := filepath.Join(workDir, "logs")
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		return err
	}

	mailer.logFile, err = os.OpenFile(
		filepath.Join(logDir, "email.log"),
		os.O_TRUNC|os.O_CREATE|os.O_RDWR,
		0644,
	)
	if err != nil {
		return err
	}

	mailer.logger = log.New(mailer.logFile, "", log.Ldate|log.Ltime)

	return nil
}

func DeinitMailer() {
	if mailer.logFile == nil {
		return
	}

	if err := mailer.logFile.Close(); err != nil {
		panic("Can't close mail file")
	}
}