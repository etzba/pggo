package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Log struct {
}

func New() *Log {
	return &Log{}
}

func (l *Log) Info(msg string) {
	fmt.Println("INFO", time.Now().UTC(), msg)
}

func (l *Log) Error(msg string, err error) {
	fmt.Println("ERROR", time.Now().UTC(), msg, err.Error())
}
