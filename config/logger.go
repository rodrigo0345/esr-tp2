package config

import (
	"fmt"
)

type Logger struct {
  level int
}

func NewLogger(level int) *Logger {
  return &Logger{
    level: level,
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
      fmt.Printf("\x1b[31m%s\x1b[0m", message)
    case 1:
      fmt.Printf("\x1b[32m%s\x1b[0m", message)
    case 2:
      fmt.Printf("\x1b[33m%s\x1b[0m", message)
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
