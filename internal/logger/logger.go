package logger

import (
	"bytes"
	"context"
	"fmt"
	"homework/internal/model"
	"text/tabwriter"
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

func (l *Logger) PrintPoints(points []model.PickUpPoint) {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 2, ' ', 0)
	fmt.Fprintf(
		w,
		"%s\t%s\t%s\t%s\n",
		"Id",
		"Name",
		"Address",
		"Contact")
	for _, point := range points {
		fmt.Fprint(w, point)
	}
	w.Flush()
	l.Log(buf.String())
}
