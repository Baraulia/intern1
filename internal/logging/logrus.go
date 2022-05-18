package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type Logger struct {
	*logrus.Entry
}

func GetLoggerLogrus() *Logger {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		FullTimestamp: true,
		ForceColors:   true,
	}
	err := os.MkdirAll("logs", 0777)
	if err != nil {
		panic(err)
	}
	allFile, err := os.OpenFile("logs/all.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	l.SetOutput(allFile)
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.TraceLevel)
	return &Logger{logrus.NewEntry(l)}
}
