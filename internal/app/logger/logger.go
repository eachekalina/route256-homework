package logger

import (
	"context"
	"fmt"
)

type Logger struct {
	out chan string
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Run(ctx context.Context) error {
	l.out = make(chan string, 128)
	defer close(l.out)
	for {
		select {
		case s := <-l.out:
			fmt.Println(s)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (l *Logger) Log(format string, a ...any) {
	l.out <- fmt.Sprintf(format, a...)
}
