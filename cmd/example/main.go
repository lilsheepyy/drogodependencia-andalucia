package main

import (
	"log"
	"os"
	"time"

	"github.com/drugprofile/drugprofile/internal/models"
	"github.com/drugprofile/drugprofile/internal/service"
)

func main() {
	pdfService := service.NewPDFService()

	// Crear perfil de ejemplo
	p := &models.Perfil{
		ID:                "7a1a9d18-83d7-4fd4-8998-4ad1a4158b93",
		NombreCompleto:    "Juan Pérez García",
		DNI:               "12345678Z",
		FechaNacimiento:   time.Date(1985, 5, 15, 0, 0, 0, 0, time.UTC),
		Direccion:         "Calle Falsa 123, Madrid",
		Telefono:          "+34 611 222 333",
		Email:             "juan.perez@ejemplo.com",
		EstadoCivil:       "Divorciado/a",
		SituacionLaboral:  "Desempleado",
		NivelEstudios:     "Bachillerato",
		FechaIngreso:      time.Now().AddDate(0, 0, -10),
		DerivadoPor:       "Servicios Sociales - Distrito Centro",
		MotivoIngreso:     "Dependencia severa a opiáceos y pérdida de red de apoyo familiar.",
		EdadInicioConsumo: 16,
		SustanciaPrincipal: "Heroína",
		FrecuenciaConsumo:  "Diario",
		ViaAdministracion:  "Inyectada",
		Sustancias: []models.Sustancia{
			{Nombre: "Alcohol"},
			{Nombre: "Benzodiacepinas"},
			{Nombre: "Tabaco"},
		},
		AntecedentesMedicos:      "Hepatitis C (tratada), anemia ferropénica crónica.",
		AntecedentesPsicologicos: "Trastorno de la personalidad límite diagnosticado en 2018.",
		Alergias:                 "Penicilina, polen, frutos secos.",
		MedicacionActual:         "Metadona 60mg/día, Diazepam 10mg noche.",
		TratamientosAnteriores:   "Ingreso en comunidad terapéutica en 2020 (abandono a los 3 meses).",
		SituacionVivienda:      "Albergue municipal",
		ProblemasLegales:       "Causa pendiente por hurto menor (juicio en septiembre).",
		ContactoEmergenciaNom:  "María García",
		ContactoEmergenciaTel:  "600 999 888",
		ParentescoEmergencia:   "Madre",
		Observaciones:          "El paciente muestra motivación para iniciar el tratamiento sustitutivo. Se requiere seguimiento estrecho por parte de trabajo social para regularizar situación de vivienda.",
	}

	pdfBytes, err := pdfService.GeneratePerfilPDF(p)
	if err != nil {
		log.Fatal("Error generando PDF de ejemplo:", err)
	}

	err = os.WriteFile("ejemplo_expediente.pdf", pdfBytes, 0644)
	if err != nil {
		log.Fatal("Error guardando el archivo PDF:", err)
	}

	log.Println("PDF de ejemplo generado con éxito: ejemplo_expediente.pdf")
}
