package reqlog

import (
	"encoding/json"
	"homework/internal/app/kafka"
	"homework/internal/app/logger"
)

func LogHandler(log logger.Logger) kafka.MessageHandler {
	return func(bytes []byte) {
		var msg Message
		err := json.Unmarshal(bytes, &msg)
		if err != nil {
			log.Log("%v\n", err)
			return
		}

		log.Log("%s, Got request method %s, path %s, params %v, headers %v, body %s\n", msg.Timestamp, msg.Method, msg.Path, msg.Params, msg.Headers, string(msg.Body))
	}
}
