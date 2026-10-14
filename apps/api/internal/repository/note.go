package repository

import (
	"database/sql"
	"fmt"

	"github.com/HiogoPariz/notez/apps/api/internal/dto"
	_ "github.com/lib/pq"
)

type NoteRepository struct {
	DB *sql.DB
}

type INoteRepository interface {
	CreateNote(*dto.NoteDTO) (int, error)
	DeleteNote(int) error
	UpdateNote(*dto.NoteDTO) error
	GetNoteByID(int) (*dto.NoteDTO, error)
	GetNotes() ([]*dto.NoteDTO, error)
}

func CreateNoteRepository(db *sql.DB) INoteRepository {
	return &NoteRepository{DB: db}
}

// CreateNote insere os metadados da nota e retorna o id gerado.
func (repo *NoteRepository) CreateNote(note *dto.NoteDTO) (int, error) {
	query := `
	INSERT INTO note (title, file_name, active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	var id int
	if err := repo.DB.QueryRow(
		query,
		note.Title,
		note.FileName,
		note.Active,
		note.CreatedAt,
		note.UpdatedAt,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (repo *NoteRepository) UpdateNote(note *dto.NoteDTO) error {
	query := `
	UPDATE note
	SET title = $1,
	    file_name = $2,
	    active = $3,
	    created_at = $4,
	    updated_at = $5
	WHERE id = $6
	`

	_, err := repo.DB.Exec(
		query,
		note.Title,
		note.FileName,
		note.Active,
		note.CreatedAt,
		note.UpdatedAt,
		note.ID,
	)

	return err
}

func (repo *NoteRepository) DeleteNote(id int) error {
	_, err := repo.DB.Exec(`UPDATE note
	 SET active = false
	 WHERE id = $1
	`, id)

	return err
}

func (repo *NoteRepository) GetNoteByID(id int) (*dto.NoteDTO, error) {
	rows, err := repo.DB.Query("SELECT * FROM note n WHERE n.id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoNote(rows)
	}
	return nil, fmt.Errorf("note.%d not found", id)
}

func (repo *NoteRepository) GetNotes() ([]*dto.NoteDTO, error) {
	rows, err := repo.DB.Query("SELECT * FROM note n WHERE n.active is true")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*dto.NoteDTO{}

	for rows.Next() {
		note, err := scanIntoNote(rows)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	return notes, nil
}

func scanIntoNote(rows *sql.Rows) (*dto.NoteDTO, error) {
	note := dto.NoteDTO{}
	err := rows.Scan(
		&note.ID,
		&note.Title,
		&note.FileName,
		&note.Active,
		&note.CreatedAt,
		&note.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &note, err
}
