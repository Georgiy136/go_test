package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"myapp/internal/models"
	"myapp/internal/sevice/nats"
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
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		rw := &responseWriter{Body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		api := c.FullPath()
		statusCode := rw.Status()
		responseBody := rw.Body.Bytes()

		log := models.Log{
			Dt:           time.Now(),
			Api:          &api,
			ServiceName:  "test_service",
			Request:      bodyBytes,
			Response:     responseBody,
			ResponseCode: &statusCode,
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
