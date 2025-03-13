package Services

import (
	"TestApp/Models"
	"TestApp/Notes/Repositories"
	userRepo "TestApp/Users/Repositories"
	"TestApp/Utils"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
)

type NotesService struct {
	Repo     Repositories.NotesRepository
	UserRepo *userRepo.UserRepository
}

type NoteService interface {
	GetAllNotes() ([]Models.Notes, error)
	ExportToCSV() ([]byte, error)
	ExportToExcel() ([]byte, error)
}

func NewNotesService(repo Repositories.NotesRepository, userRepo *userRepo.UserRepository) *NotesService {
	return &NotesService{Repo: repo, UserRepo: userRepo}
}

func (service *NotesService) CreateNote(note *Models.Notes) error {
	if note == nil {
		return errors.New("Data can't be null")
	}
	if note.NoteTitle == "" || note.NoteMsg == "" {
		return errors.New("title and content cannot be empty")
	}

	if service.UserRepo.FindUserById(note.UserID) {
		return fmt.Errorf("User not found")
	}

	return service.Repo.CreateNote(note)
}

func (service *NotesService) GetAllNotes(userId uint64, query string) ([]Models.Notes, error) {
	return service.Repo.GetAllNotes(userId, query)
}

func (service *NotesService) GetNoteById(id uint64) (Models.Notes, error) {
	return service.Repo.GetNoteById(id)
}

func (service *NotesService) UpdateNote(id uint64, note *Models.Notes) (*Models.Notes, error) {
	updatedNote, err := service.Repo.UpdateNote(id, note)
	if err != nil {
		return nil, err
	}
	return updatedNote, nil
}

func (service *NotesService) DeleteNote(id uint64) error {
	return service.Repo.DeleteNote(id)
}

func (service *NotesService) ExportToCSV(userId uint64) ([]byte, error) {
	notes, err := service.Repo.GetAllNotes(userId, "")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"Title", "Note"})

	for _, note := range notes {
		writer.Write([]string{
			note.NoteTitle,
			note.NoteMsg,
		})
	}
	writer.Flush()
	return buf.Bytes(), nil
}

func (service *NotesService) ExportToExcel(userId uint64) ([]byte, error) {
	notes, err := service.Repo.GetAllNotes(userId, "")
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	sheetName := "exported Notes"
	file.SetSheetName(file.GetSheetName(0), sheetName)

	headers := []string{"Title", "Note"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string('A'+i))
		file.SetCellValue(sheetName, header, cell)
	}

	for idx, note := range notes {
		row := idx + 2
		file.SetCellValue(sheetName, fmt.Sprintf("B%d", row), note.NoteTitle)
		file.SetCellValue(sheetName, fmt.Sprintf("C%d", row), note.NoteMsg)
	}

	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (service *NotesService) ExportToPDF(ctx echo.Context, id uint64) error {
	note, err := service.Repo.GetNoteById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Note not found"})
	}

	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithRightMargin(15).
		WithBottomMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(cfg)

	// Add content to PDF
	Utils.ToPDF(m, note.NoteTitle, note.NoteMsg)

	document, err := m.Generate()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate PDF"})
	}

	filePath := fmt.Sprintf("pdf/note_%d.pdf", note.ID)
	err = document.Save(filePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"msg":    "Failed to save PDF",
			"error:": err.Error(),
		})
	}
	defer os.Remove(filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open PDF"})
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get file size"})
	}

	// Set response headers
	ctx.Response().Header().Set("Content-Type", "application/pdf")
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=note_%d.pdf", note.ID))
	ctx.Response().Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))

	_, err = io.Copy(ctx.Response().Writer, file)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to stream PDF"})
	}
	return nil
}

func (service *NotesService) ImportFromExcel(ctx echo.Context, userId uint64) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "Unable to get the file",
		})
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error while opening the file",
		})
	}
	defer src.Close()

	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error reading the Excel file",
		})
	}

	sheet := xlsx.GetSheetName(0)

	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"message": "Error reading rows from the Excel file",
		})
	}

	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) < 2 {
			continue
		}
		note := &Models.Notes{
			UserID:    uint(userId),
			NoteTitle: row[0],
			NoteMsg:   row[1],
		}

		err = service.Repo.CreateNote(note)
		if err != nil {
			log.Printf("Error creating note from row: %v", err)
		}
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"message": "Notes imported successfully",
	})
}
