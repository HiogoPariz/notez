package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/HiogoPariz/notez/apps/api/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB) {
	router := gin.Default()
	noteService := createNoteService(db)

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.Get("WEB_URL", "http://localhost:5173")},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	router.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	router.GET("/note", noteService.GetNotes)
	router.GET("/note/:id", noteService.GetNoteByID)
	router.POST("/note", noteService.CreateNote)

	if err := router.Run(":" + config.Get("API_PORT", "3000")); err != nil {
		panic(err)
	}
}
