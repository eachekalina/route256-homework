package logger

import (
	"bytes"
	"fmt"
	"homework/internal/model"
	"text/tabwriter"
)

type Logger struct {
	out chan string
}

func NewLogger() *Logger {
	return &Logger{out: make(chan string, 128)}
}

func (l *Logger) Run() {
	go func() {
		for {
			s, ok := <-l.out
			if !ok {
				return
			}
			fmt.Println(s)
		}
	}()
}

func (l *Logger) Close() {
	close(l.out)
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
