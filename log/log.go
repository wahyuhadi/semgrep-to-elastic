package log

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func New(appName, logDir string) (*logrus.Logger, *os.File) {
	var file *os.File

	filename := fmt.Sprintf("%s.log", appName)

	dir := logDir
	if strings.HasSuffix(dir, "/") {
		strings.TrimSuffix(dir, "/")
	}

	instance := logrus.New()
	if f, err := os.OpenFile(dir+filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755); err == nil {
		o := io.MultiWriter(os.Stdout, f)
		instance.SetOutput(o)
	} else {
		instance.Warnln("Error create logs")
	}

	return instance, file
}
