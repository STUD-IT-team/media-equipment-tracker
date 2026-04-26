package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(NewTextFormatter(true))
	logrus.SetOutput(os.Stdout)
}
