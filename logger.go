package kurin

import (
	"log"
	"os"
)

type (
	Logger interface {
		Debug(args ...interface{})
		Info(args ...interface{})
		Warn(args ...interface{})
		Error(args ...interface{})
		Fatal(args ...interface{})
		Panic(args ...interface{})
	}

	defaultLogger struct {
		stdout *log.Logger
		stderr *log.Logger
	}
)

func newDefaultLogger() Logger {
	return &defaultLogger{
		stdout: log.New(os.Stdout, "", log.LstdFlags),
		stderr: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (logger *defaultLogger) Debug(args ...interface{}) {
	logger.stdout.Println(args...)
}

func (logger *defaultLogger) Info(args ...interface{}) {
	logger.stdout.Println(args...)
}

func (logger *defaultLogger) Warn(args ...interface{}) {
	logger.stdout.Println(args...)
}

func (logger *defaultLogger) Error(args ...interface{}) {
	logger.stderr.Println(args...)
}

func (logger *defaultLogger) Fatal(args ...interface{}) {
	logger.stderr.Fatalln(args...)
}

func (logger *defaultLogger) Panic(args ...interface{}) {
	logger.stderr.Panicln(args...)
}
