package log

import (
	logs "github.com/sirupsen/logrus"
	"os"
)

// Initialize log handle
func LogInit() {
	logs.SetOutput(os.Stdout)
	logs.SetFormatter(&logs.TextFormatter{
		//	DisableColors: true,
		FullTimestamp: true,
	})
}
func Fatal(format string, args ...interface{}) {
	if len(args) > 0 {
		logs.Fatalf(format, args)
	} else {
		logs.Fatal(format)
	}
}

func Info(format string, args ...interface{}) {
	if len(args) > 0 {
		logs.Infof(format, args)
	} else {
		logs.Info(format)
	}
}
