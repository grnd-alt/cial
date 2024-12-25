package services

import (
	"backendsetup/m/db/sql/dbgen"
	"context"
)

type NotesService struct {
	query *dbgen.Queries
}

func InitNotesService(queries *dbgen.Queries) *NotesService {
	return &NotesService{
		query: queries,
	}
}

func (n *NotesService) CreateNote(createdBy string, title string, content string) (*dbgen.Note, error) {
	note, err := n.query.CreateNote(context.Background(), dbgen.CreateNoteParams{CreatedBy: createdBy, Title: title, Content: content})
	return &note, err
}


func (n *NotesService) GetNotes(createdBy string) ([]dbgen.Note, error) {
	notes, err := n.query.GetNotesByUser(context.Background(), createdBy)
	return notes, err
}
