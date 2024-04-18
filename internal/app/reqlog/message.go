package reqlog

import (
	"net/http"
	"time"
)

type Message struct {
	Timestamp time.Time
	Method    string
	Path      string
	Headers   http.Header
	Params    map[string]string
	Body      string
}
