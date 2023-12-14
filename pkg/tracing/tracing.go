package tracing

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Trace struct {
	RequestID     string
	RequestURL    string
	RequestMethod string
}

func NewTrace(c *gin.Context) *Trace {
	return &Trace{
		RequestID:     uuid.NewString(),
		RequestURL:    c.Request.RequestURI,
		RequestMethod: c.Request.Method,
	}
}

func TraceMiddleware(c *gin.Context) {
	trace := NewTrace(c)

	logrus.WithFields(logrus.Fields{
		"request_id":     trace.RequestID,
		"request_url":    trace.RequestURL,
		"request_method": trace.RequestMethod,
	}).Info("Request received")

	c.Set("requestID", trace)
	c.Next()
}
