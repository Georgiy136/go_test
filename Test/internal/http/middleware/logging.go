package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

type responseWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.Body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		rw := &responseWriter{Body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = rw

		c.Next()

		statusCode := rw.Status()
		responseBody := rw.Body.String()

		// отправляем лог в Nats
		// ...

		fmt.Printf("Request Body: %s\n", string(bodyBytes))
		fmt.Printf("Response Status Code: %d\n", statusCode)
		fmt.Printf("Response Body: %s\n", responseBody)
	}
}
