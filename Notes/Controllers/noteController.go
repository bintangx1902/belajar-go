package Controllers

import (
	"TestApp/Models"
	"TestApp/Notes/Services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type NoteController struct {
	Service *Services.NotesService
}

func NewNotesController(service *Services.NotesService) *NoteController {
	return &NoteController{Service: service}
}

func (controller *NoteController) CreateNote(ctx echo.Context) error {
	var note Models.Notes

	if err := ctx.Bind(&note); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if (note.NoteTitle == "" || note.NoteMsg == "") && note.UserID == 0 {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "NoteTitle and NoteMsg cannot be empty"})
	}

	if controller.Service.UserRepo.FindUserById(note.UserID) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "User not found"})
	}

	err := controller.Service.CreateNote(&note)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, map[string]string{"created": "ok"})
}

func (controller *NoteController) GetAllNotes(ctx echo.Context) error {
	query := ctx.QueryParam("q")
	userId := ctx.Get("user_id").(uint64)
	notes, err := controller.Service.GetAllNotes(userId, query)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.JSON(http.StatusOK, notes)
}

func (controller *NoteController) GetNoteById(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}
	note, err := controller.Service.GetNoteById(id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, note)
}

func (controller *NoteController) UpdateNote(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}
	var note Models.Notes
	if err := ctx.Bind(&note); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	updateNote, err := controller.Service.UpdateNote(id, &note)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, updateNote)
}

func (controller *NoteController) DeleteNoteById(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid ID parameter",
		})
	}
	if err := controller.Service.DeleteNote(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (controller *NoteController) DownloadCSV(ctx echo.Context) error {
	userId := ctx.Get("user_id").(uint64)
	data, err := controller.Service.ExportToCSV(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.Blob(http.StatusOK, "text/csv", data)
}

func (controller *NoteController) DownloadExcel(ctx echo.Context) error {
	userId := ctx.Get("user_id").(uint64)
	data, err := controller.Service.ExportToExcel(userId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (controller *NoteController) DownloadPDF(ctx echo.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID parameter"})
	}

	err = controller.Service.ExportToPDF(ctx, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return nil
}

func (controller *NoteController) ImportFromExcel(ctx echo.Context) error {
	id := ctx.Get("user_id").(uint64)
	errs := controller.Service.ImportFromExcel(ctx, id)
	if errs != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": errs.Error()})
	}
	return ctx.JSON(http.StatusOK, map[string]string{"message": "Successfully Imported"})
}
