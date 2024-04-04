//go:generate mockgen -source=./logger.go -destination=../mocks/logger.go -package=mocks

package logger

import (
	"context"
	"fmt"
)

type Logger interface {
	Log(format string, a ...any)
}

type ThreadLogger struct {
	out chan string
}

func NewLogger() *ThreadLogger {
	return &ThreadLogger{}
}

func (l *ThreadLogger) Run(ctx context.Context) error {
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

func (l *ThreadLogger) Log(format string, a ...any) {
	l.out <- fmt.Sprintf(format, a...)
}
