package utils

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LogFormatter struct {
}

func (s *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("0102-150405.000")
	var file string
	var len int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		len = entry.Caller.Line
	}
	//fmt.Println(entry.Data)
	msg := fmt.Sprintf("%s [%s:%d][GOID:%d][%s] %s\n", timestamp, file, len, getGID(), strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}

func getGID() uint64 {

	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

type logFileWriter struct {
	file     *os.File
	logPath  string
	fileDate string // judge the date to switch the directory
	appName  string
	encoding string
}

func (p *logFileWriter) Write(data []byte) (n int, err error) {
	if p == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if p.file == nil {
		return 0, errors.New("file not opened")
	}

	//Determine whether to switch dates
	fileDate := time.Now().Format("20060102")
	if p.fileDate != fileDate {
		p.file.Close()
		err = os.MkdirAll(fmt.Sprintf("%s/%s", p.logPath, fileDate), os.ModePerm)
		if err != nil {
			return 0, err
		}
		filename := fmt.Sprintf("%s/%s/%s-%s.log", p.logPath, fileDate, p.appName, fileDate)

		p.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		if err != nil {
			return 0, err
		}

	}
	//if p.encoding != "" {
	//	dataToEncode := ConvertStringToByte(string(data), p.encoding)
	//	n, e := p.file.Write(dataToEncode)
	//	return n, e
	//}

	n, e := p.file.Write(data)
	return n, e

}

// InitLog Initialization log
func InitLog(logPath string, appName string, encoding string) {
	fileDate := time.Now().Format("20060102")
	//Create directory
	err := os.MkdirAll(fmt.Sprintf("%s/%s", logPath, fileDate), os.ModePerm)
	if err != nil {
		log.Error(err)
		return
	}

	filename := fmt.Sprintf("%s/%s/%s-%s.log", logPath, fileDate, appName, fileDate)
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		log.Error(err)
		return
	}

	fileWriter := logFileWriter{file, logPath, fileDate, appName, encoding}
	log.SetOutput(&fileWriter)

	log.SetReportCaller(true)
	log.SetFormatter(new(LogFormatter))
}
