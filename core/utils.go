package core

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WithDBCallback func(d *sql.DB, c *gin.Context) error
type WithDBResult func(c *gin.Context) error

func WithDB(d *sql.DB, callback WithDBCallback) WithDBResult {
	return func(c *gin.Context) error {
		return callback(d, c)
	}
}

func WrapError(wrapped WithDBResult) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := wrapped(c); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}
}
