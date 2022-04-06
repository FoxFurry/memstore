/*
Package httperr

Responsible for creating easy way of passing both error message and status code to handler
*/
package httperr

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// HttpErr implements error interface and contains status code and message
//
// HttpErr can be used with Handle method to properly write it to gin context
type HttpErr struct {
	statusCode int
	message    string
}

// Error implements error interface for HttpErr
func (e HttpErr) Error() string {
	return e.message
}

// Handle writes given error to gin context
// If giver error is HttpErr, it will extract status code and message and use them
// Otherwise it will use 500 status code and marshal error
func Handle(c *gin.Context, err error) {
	if hErr, ok := err.(HttpErr); ok { // If it is httperr - get status code from it
		c.JSON(hErr.statusCode, gin.H{
			"error": hErr.message,
		})
	} else { // Otherwise use internal error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}

// New creates HttpErr from given message and status code
func New(msg string, code int) error {
	return HttpErr{
		message:    msg,
		statusCode: code,
	}
}

// Wrap creates a generic error by wrapping given error into message
func Wrap(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s: %v", msg, err))
}

// WrapHttp creates a HttpErr by wrapping given error into message and adding status code
func WrapHttp(err error, msg string, code int) error {
	return HttpErr{
		message:    fmt.Sprintf("%s: %v", msg, err),
		statusCode: code,
	}
}
