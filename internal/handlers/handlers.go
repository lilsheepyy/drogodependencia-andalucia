package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/drugprofile/drugprofile/internal/models"
	"github.com/drugprofile/drugprofile/internal/repository"
	"github.com/drugprofile/drugprofile/internal/service"
	"github.com/drugprofile/drugprofile/internal/views"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

type Handlers struct {
	repo *repository.Repository
	pdf  *service.PDFService
}

func NewHandlers(repo *repository.Repository, pdf *service.PDFService) *Handlers {
	return &Handlers{repo: repo, pdf: pdf}
}

func Render(c *fiber.Ctx, component templ.Component) error {
	return adaptor.HTTPHandler(templ.Handler(component))(c)
}

func (h *Handlers) Index(c *fiber.Ctx) error {
	sustancias, err := h.repo.GetSustancias("")
	if err != nil {
		// Log error but continue with empty list
		return Render(c, views.Index([]models.Sustancia{}))
	}
	return Render(c, views.Index(sustancias))
}

func (h *Handlers) SearchSustancias(c *fiber.Ctx) error {
	search := c.FormValue("search_filter")
	sustancias, err := h.repo.GetSustancias(search)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return Render(c, views.DrugList(sustancias))
}

func (h *Handlers) SearchPatient(c *fiber.Ctx) error {
	dni := c.FormValue("dni_search")
	if dni == "" {
		return Render(c, views.PatientSearchResults(nil, "Ingrese DNI o Pasaporte para buscar."))
	}

	p, err := h.repo.GetPerfilByDNI(dni)
	if err != nil {
		return Render(c, views.PatientSearchResults(nil, "No se encontró ningún paciente con ese documento."))
	}

	return Render(c, views.PatientSearchResults(p, ""))
}

func (h *Handlers) NewReportForm(c *fiber.Ctx) error {
	dni := c.Params("dni")
	p, _ := h.repo.GetPerfilByDNI(dni) // If not found, p will be nil
	
	sustancias, err := h.repo.GetSustancias("")
	if err != nil {
		return Render(c, views.Index(nil))
	}
	
	component := views.RegistrationForm(sustancias, p)
	
	// Si no es una petición HTMX (ej: recarga de página), envolver en Dashboard (que incluye layout y cabecera)
	if c.Get("HX-Request") == "" {
		return Render(c, views.Dashboard("Nuevo Informe", component))
	}
	
	return Render(c, component)
}

func (h *Handlers) AddManualSustancia(c *fiber.Ctx) error {
	nombre := c.FormValue("manual_name")
	if nombre == "" {
		return c.Status(400).SendString("Nombre vacío")
	}

	s, err := h.repo.CreateSustanciaManual(nombre)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return Render(c, views.ManualDrug(*s))
}


func (h *Handlers) CreatePerfil(c *fiber.Ctx) error {
	fechaNac, _ := time.Parse("02/01/2006", c.FormValue("fecha_nacimiento"))
	fechaIng, _ := time.Parse("02/01/2006", c.FormValue("fecha_ingreso"))
	edadInicio, _ := strconv.Atoi(c.FormValue("edad_inicio"))

	// Captura de IDs de sustancias (manejo robusto para múltiples valores)
	var ids []string
	
	// Intentar obtener de multipart (si se envió así)
	form, err := c.MultipartForm()
	if err == nil && form.Value != nil {
		ids = form.Value["sustancias_ids"]
	}
	
	// Si no hay nada en multipart, intentar capturar manualmente de la petición form-encoded
	if len(ids) == 0 {
		// Fiber c.BodyParser también puede funcionar con un struct, 
		// pero aquí extraemos manualmente para mayor control
		c.Context().PostArgs().VisitAll(func(key, value []byte) {
			if string(key) == "sustancias_ids" {
				ids = append(ids, string(value))
			}
		})
	}

	var sustancias []models.Sustancia
	for _, id := range ids {
		if id != "" {
			sustancias = append(sustancias, models.Sustancia{ID: id})
		}
	}

	p := &models.Perfil{
		// 1. Datos Personales
		NombreCompleto:   c.FormValue("nombre_completo"),
		DNI:              c.FormValue("dni"),
		FechaNacimiento:  fechaNac,
		Direccion:        c.FormValue("direccion"),
		Telefono:         c.FormValue("telefono"),
		Email:            c.FormValue("email"),
		EstadoCivil:      c.FormValue("estado_civil"),
		SituacionLaboral: c.FormValue("situacion_laboral"),
		NivelEstudios:    c.FormValue("nivel_estudios"),

		// 2. Historial de Ingreso
		FechaIngreso:  fechaIng,
		MotivoIngreso: c.FormValue("motivo_ingreso"),
		DerivadoPor:   c.FormValue("derivado_por"),

		// 3. Historial de Consumo
		EdadInicioConsumo:  edadInicio,
		SustanciaPrincipal: c.FormValue("sustancia_principal"),
		FrecuenciaConsumo:  c.FormValue("frecuencia"),
		ViaAdministracion:  c.FormValue("via_admin"),
		Sustancias:         sustancias,

		// 4. Salud y Tratamiento
		AntecedentesMedicos:      c.FormValue("ant_medicos"),
		AntecedentesPsicologicos: c.FormValue("ant_psicologicos"),
		Alergias:                 c.FormValue("alergias"),
		MedicacionActual:         c.FormValue("medicacion"),
		TratamientosAnteriores:   c.FormValue("trat_anteriores"),

		// 5. Situación Social y Legal
		SituacionVivienda:     c.FormValue("vivienda"),
		ProblemasLegales:      c.FormValue("prob_legales"),
		ContactoEmergenciaNom: c.FormValue("emergencia_nom"),
		ContactoEmergenciaTel: c.FormValue("emergencia_tel"),
		ParentescoEmergencia:  c.FormValue("emergencia_parentesco"),

		Observaciones: c.FormValue("observaciones"),
	}

	err = h.repo.CreatePerfil(p)
	if err != nil {
		return Render(c, views.Result("", false, "Error al guardar en base de datos: "+err.Error()))
	}

	return Render(c, views.Result(p.ID, true, "El expediente clínico ha sido guardado de forma segura."))
}

func (h *Handlers) GetPDF(c *fiber.Ctx) error {
	id := c.Params("id")
	p, err := h.repo.GetPerfil(id)
	if err != nil {
		return c.Status(404).SendString("Perfil clínico no encontrado en el sistema.")
	}

	pdfBytes, err := h.pdf.GeneratePerfilPDF(p)
	if err != nil {
		return c.Status(500).SendString("Error crítico en la generación del reporte PDF: " + err.Error())
	}

	c.Set("Content-Type", "application/pdf")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"expediente_%s.pdf\"", p.NombreCompleto))
	return c.Send(pdfBytes)
}
