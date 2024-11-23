package config

import (
	"fmt"
	"time"
)

type Logger struct {
	level int
  nodeName string
}

func NewLogger(level int, nodeName string) *Logger {
	return &Logger{
		level: level,
    nodeName: nodeName,
	}
}

func (l *Logger) SetLevel(level int) {
	l.level = level
}

func (l *Logger) Log(level int, message string) {
	if level <= l.level {
		// add colors depending on the level
		switch level {
		case 0:
			fmt.Printf("\x1b[31m[%s - %s] %s\x1b[0m\n", time.Now().Format("15:04:05"), l.nodeName, message)
		case 1:
			fmt.Printf("\x1b[32m[%s - %s] %s\x1b[0m\n", time.Now().Format("15:04:05"), l.nodeName, message)
		case 2:
			fmt.Printf("\x1b[33m[%s - %s] %s\x1b[0m\n", time.Now().Format("15:04:05"), l.nodeName, message)
		}
	}
}

func (l *Logger) Info(message string) {
	l.Log(1, message)
}

func (l *Logger) Error(message string) {
	l.Log(0, message)
}

func (l *Logger) Debug(message string) {
	l.Log(2, message)
}
