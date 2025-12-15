package logging

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// New creates a logger that writes to logs/<component>.log and returns it with a cleanup.
func New(component string) (*logrus.Entry, func(), error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, nil, err
	}
	path := filepath.Join("logs", component+".log")
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}

	logger.SetOutput(f)
	return logger.WithField("component", component), func() { _ = f.Close() }, nil
}
