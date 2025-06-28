package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/Georgiy136/go_test/go_test/internal/models"
	"github.com/Georgiy136/go_test/go_test/internal/nats"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func NewLogger(nats *nats.NatsService, channelName string) *Logger {
	return &Logger{nats: nats, channelName: channelName}
}

type Logger struct {
	nats        *nats.NatsService
	channelName string
}

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (l *Logger) LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody []byte
		if c.Request.Body != nil {
			reqBody, _ = io.ReadAll(c.Request.Body)
		}

		// Восстановление тела запроса
		c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		rw := &responseWriter{Body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		api := c.FullPath()
		statusCode := rw.Status()
		responseBody := rw.Body.String()

		log := models.Log{
			Dt:           time.Now(),
			Api:          api,
			ServiceName:  "test_service",
			Request:      string(reqBody),
			Response:     responseBody,
			ResponseCode: statusCode,
		}

		logBytes, err := json.Marshal(log)
		if err != nil {
			logrus.Errorf("loggingMiddleware error marshalling log: %v", err)
		}

		if err = l.nats.SendBatch(l.channelName, logBytes); err != nil {
			logrus.Errorf("loggingMiddleware error publish log: %v", err)
		}
	}
}
