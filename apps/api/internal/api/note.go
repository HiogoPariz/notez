package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/HiogoPariz/notez/apps/api/internal/dto"
	"github.com/HiogoPariz/notez/apps/api/internal/integration"
	"github.com/HiogoPariz/notez/apps/api/internal/repository"
	"github.com/gin-gonic/gin"
)

type NoteService struct {
	DB *sql.DB
}

func createNoteService(db *sql.DB) *NoteService {
	return &NoteService{DB: db}
}

func (service NoteService) GetNotes(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)
	notes_dto, err := repo.GetNotes()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	notes := []*dto.NoteObject{}

	for _, note_dto := range notes_dto {
		note_res, err := note_dto.EntityToResponse()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		content, err := integration.CreateFileIntegration(note_dto).GetFileContent()
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		note_res.Content = content

		notes = append(notes, note_res)
	}

	ctx.JSON(http.StatusOK, notes)
}

func (service NoteService) GetNoteByID(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	note_dto, err := repo.GetNoteByID(id)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, err)
		return
	}

	note, err := note_dto.EntityToResponse()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	content, err := integration.CreateFileIntegration(note_dto).GetFileContent()
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	note.Content = content

	ctx.JSON(http.StatusOK, note)
}

func (service NoteService) CreateNote(ctx *gin.Context) {
	repo := repository.CreateNoteRepository(service.DB)

	body := &dto.NoteRequest{}
	if err := ctx.Bind(body); err != nil {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	note_dto := body.RequestToDTO("")

	// Grava os metadados primeiro para obter o id gerado pelo banco.
	id, err := repo.CreateNote(note_dto)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// O id e unico, entao serve como base para o nome do arquivo de conteudo.
	file_name := fmt.Sprintf("note-%d", id)
	note_dto.ID = id
	note_dto.FileName = file_name

	if err := repo.UpdateNote(note_dto); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Manda o conteudo pro servico de files.
	if err := integration.CreateFileIntegration(note_dto).CreateFileContent(body.Content, file_name); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, note_dto)
}
