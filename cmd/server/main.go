package main

import (
	"log"

	"github.com/drugprofile/drugprofile/internal/handlers"
	"github.com/drugprofile/drugprofile/internal/repository"
	"github.com/drugprofile/drugprofile/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	// Inicializar base de datos
	repo, err := repository.NewRepository("drugprofile.db")
	if err != nil {
		log.Fatal("Error inicializando base de datos:", err)
	}

	// Inicializar servicios
	pdfService := service.NewPDFService()
	store := session.New()

	// Inicializar handlers
	h := handlers.NewHandlers(repo, pdfService, store)

	// Configurar Fiber
	app := fiber.New(fiber.Config{
		AppName: "Gestión de Drogodependencia",
	})

	app.Use(logger.New())

	// Rutas Públicas (Auth)
	app.Get("/login", h.LoginForm)
	app.Post("/login", h.Login)
	app.Get("/logout", h.Logout)

	// Rutas Protegidas
	protected := app.Group("/", h.AuthMiddleware)
	
	protected.Get("/cambiar-password", h.PasswordChangeForm)
	protected.Post("/cambiar-password", h.UpdatePassword)

	protected.Get("/", h.Index)
	protected.Post("/paciente/search", h.SearchPatient)
	protected.Get("/perfil/nuevo", h.NewReportForm)
	protected.Get("/perfil/nuevo/:dni", h.NewReportForm)
	protected.Post("/sustancias/search", h.SearchSustancias)
	protected.Post("/sustancias/manual", h.AddManualSustancia)
	protected.Post("/perfil", h.CreatePerfil)
	protected.Get("/pdf/:id", h.GetPDF)

	// Rutas de Administración
	admin := protected.Group("/admin", h.AdminMiddleware)
	admin.Get("/", h.AdminPanel)
	admin.Post("/usuarios", h.CreateUser)
	admin.Delete("/usuarios/:id", h.DeleteUser)

	// Servir archivos estáticos
	app.Static("/static", "./static")

	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
