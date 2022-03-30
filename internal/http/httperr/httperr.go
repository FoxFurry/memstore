package httperr

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleError(c *gin.Context, err error, statusCode int) {
	c.JSON(statusCode, gin.H{
		"error": err,
	})
}

func InternalError(c *gin.Context, err error) {
	handleError(c, err, http.StatusInternalServerError)
}

func BadRequest(c *gin.Context, err error) {
	handleError(c, err, http.StatusBadRequest)
}

func Wrap(err error, msg string) error {
	return errors.New(fmt.Sprintf("%s: %v", msg, err))
}
