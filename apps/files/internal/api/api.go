package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/HiogoPariz/notez/apps/files/internal/storage"
	"github.com/gin-gonic/gin"
)

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(corsMiddleware())

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.GET("/:path", func(c *gin.Context) {
		path := c.Params.ByName("path")
		var file map[string]interface{}

		content, err := storage.ReadFile(path)
		if err != nil {
			abortWithStorageError(c, err)
			return
		}

		if err := json.Unmarshal([]byte(content), &file); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, &file)
	})

	router.POST("/:path", func(c *gin.Context) {
		path := c.Params.ByName("path")

		content, err := c.GetRawData()
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if err := storage.WriteFile(path, string(content)); err != nil {
			abortWithStorageError(c, err)
			return
		}

		result := Message{Type: "success", Message: "created " + path, Code: 0}

		c.JSON(http.StatusOK, result)
	})

	router.DELETE("/:path", func(c *gin.Context) {
		path := c.Params.ByName("path")

		if err := storage.DeleteFile(path); err != nil {
			abortWithStorageError(c, err)
			return
		}

		result := Message{Type: "success", Message: "deleted " + path, Code: 0}

		c.JSON(http.StatusOK, result)
	})

	return router
}

// abortWithStorageError maps storage errors to proper HTTP status codes.
func abortWithStorageError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, storage.ErrInvalidPath):
		c.AbortWithError(http.StatusBadRequest, err)
	case errors.Is(err, os.ErrNotExist):
		c.AbortWithError(http.StatusNotFound, err)
	default:
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}
