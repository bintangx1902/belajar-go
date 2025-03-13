package Repositories

import (
	"TestApp/Models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type NotesRepository struct {
	DB *gorm.DB
}

type NoteRepository interface {
	GetAllNotes(query string) ([]Models.Notes, error)
}

func NewNotesRepository(db *gorm.DB) NotesRepository {
	return NotesRepository{DB: db}
}

func (r *NotesRepository) CreateNote(note *Models.Notes) error {
	return r.DB.Create(note).Error
}

func (r *NotesRepository) GetAllNotes(userId uint64, queryFilter string) ([]Models.Notes, error) {
	var notes []Models.Notes
	query := r.DB.Model(&Models.Notes{}).Where("user_id = ?", userId)

	if queryFilter != "" {
		query = query.Where("note_title LIKE ? or note LIKE ?", "%"+queryFilter+"%", "%"+queryFilter+"%")
	}
	err := query.Preload("Users").Find(&notes).Error
	return notes, err
}

func (r *NotesRepository) GetNoteById(id uint64) (Models.Notes, error) {
	var note Models.Notes
	if id == 0 {
		return note, fmt.Errorf("ID Can't be zero or Noting")
	}
	err := r.DB.Model(&Models.Notes{}).
		Preload("Users").
		Where("id = ?", id).
		First(&note, id).
		Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return note, fmt.Errorf("note With ID = %d Not Found", id)
		}
		return note, err
	}
	return note, nil
}

func (r *NotesRepository) UpdateNote(id uint64, note *Models.Notes) (*Models.Notes, error) {
	var exNote Models.Notes
	err := r.DB.Model(&exNote).First(&exNote, id).Error
	if err != nil {
		return nil, fmt.Errorf("note not found")
	}
	exNote.NoteTitle = note.NoteTitle
	exNote.NoteMsg = note.NoteMsg

	if err := r.DB.Save(&exNote).Error; err != nil {
		return nil, fmt.Errorf("could not update the note because : %v", err)
	}
	return &exNote, nil
}

func (r *NotesRepository) DeleteNote(id uint64) error {
	var note Models.Notes
	err := r.DB.Model(&Models.Notes{}).First(&note, id).Error
	if err != nil {
		return fmt.Errorf("note not found")
	}
	err = r.DB.Delete(&note).Error
	if err != nil {
		return fmt.Errorf("could not delete note: %v", err)
	}
	return nil
}
