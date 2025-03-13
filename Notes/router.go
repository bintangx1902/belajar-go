package Notes

import (
	"TestApp/Middlewares"
	"TestApp/Notes/Controllers"
	"TestApp/Notes/Repositories"
	"TestApp/Notes/Services"
	userRepository "TestApp/Users/Repositories"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func NotesRouter(route *echo.Group, db *gorm.DB) {
	repo := Repositories.NewNotesRepository(db)
	userRepo := userRepository.NewUsersRepository(db)
	service := Services.NewNotesService(repo, userRepo)
	controller := Controllers.NewNotesController(service)

	notesGroup := route.Group("/notes")

	// GET
	notesGroup.GET("/", controller.GetAllNotes, Middlewares.JWTMiddleware)
	notesGroup.GET("/:id", controller.GetNoteById, Middlewares.JWTMiddleware)
	notesGroup.GET("/:id/download/pdf", controller.DownloadPDF)
	notesGroup.GET("/download/csv", controller.DownloadCSV, Middlewares.JWTMiddleware)
	notesGroup.GET("/download/excel", controller.DownloadExcel, Middlewares.JWTMiddleware)

	// POST
	notesGroup.POST("/", controller.CreateNote, Middlewares.JWTMiddleware)
	notesGroup.POST("/upload/excel", controller.ImportFromExcel, Middlewares.JWTMiddleware)

	// UPDATE
	notesGroup.PUT("/:id", controller.UpdateNote, Middlewares.JWTMiddleware)

	// DELETE
	notesGroup.DELETE("/:id", controller.DeleteNoteById, Middlewares.JWTMiddleware)

}
