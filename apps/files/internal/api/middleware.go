package api

import (
	"time"

	"github.com/HiogoPariz/notez/apps/files/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{config.Get("WEB_URL", "http://localhost:5173")},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		MaxAge:       12 * time.Hour,
	})
}
