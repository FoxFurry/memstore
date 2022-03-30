package httperr

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpErr struct {
	statusCode int
	msg        string
}

func (e HttpErr) Error() string {
	return e.msg
}

func handleError(c *gin.Context, err error) {
	if httperr, ok := err.(HttpErr); ok {
		c.JSON(httperr.statusCode, gin.H{
			"error": httperr.msg,
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
	}
}

func InternalError(c *gin.Context, err error) {
	handleError(c, err)
}

func BadRequest(c *gin.Context, err error) {
	handleError(c, err)
}

func Wrap(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s: %v", msg, err))
}
