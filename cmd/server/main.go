package main

import (
	"log"

	"github.com/drugprofile/drugprofile/internal/handlers"
	"github.com/drugprofile/drugprofile/internal/repository"
	"github.com/drugprofile/drugprofile/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Inicializar base de datos
	repo, err := repository.NewRepository("drugprofile.db")
	if err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}

	// Inicializar servicios
	pdfService := service.NewPDFService()

	// Inicializar handlers
	h := handlers.NewHandlers(repo, pdfService)

	// Configurar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Gestión de Drogodependencia",
	})

	app.Use(logger.New())

	// Rutas
	app.Get("/", h.Index)
	app.Post("/paciente/search", h.SearchPatient)
	app.Get("/perfil/nuevo", h.NewReportForm)
	app.Get("/perfil/nuevo/:dni", h.NewReportForm)
	app.Post("/sustancias/search", h.SearchSustancias)
	app.Post("/sustancias/manual", h.AddManualSustancia)
	app.Post("/perfil", h.CreatePerfil)
	app.Get("/pdf/:id", h.GetPDF)

	// Servir archivos estáticos
	app.Static("/static", "./static")

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
