package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Sustancia representa una droga o sustancia de consumo
type Sustancia struct {
	ID       string `gorm:"primaryKey" json:"id"`
	Nombre   string `gorm:"uniqueIndex;not null" json:"nombre"`
	EsManual bool   `gorm:"default:false" json:"es_manual"` // Flag para distinguir de la lista fija
}

func (s *Sustancia) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}

// Perfil representa los datos detallados de una persona en drogodependencia
type Perfil struct {
	ID                string      `gorm:"primaryKey" json:"id"`
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// 1. Datos Personales
	NombreCompleto    string      `gorm:"not null" json:"nombre_completo"`
	DNI               string      `json:"dni"`
	FechaNacimiento   time.Time   `json:"fecha_nacimiento"`
	Direccion         string      `json:"direccion"`
	Telefono          string      `json:"telefono"`
	Email             string      `json:"email"`
	EstadoCivil       string      `json:"estado_civil"`       // Soltero/a, Casado/a, Divorciado/a, Viudo/a, Pareja de hecho
	SituacionLaboral  string      `json:"situacion_laboral"`  // Empleado, Desempleado, Estudiante, Jubilado, Incapacitado
	NivelEstudios     string      `json:"nivel_estudios"`

	// 2. Historial de Ingreso
	FechaIngreso      time.Time   `json:"fecha_ingreso"`
	MotivoIngreso     string      `json:"motivo_ingreso"`
	DerivadoPor       string      `json:"derivado_por"`       // Quien lo envía (Médico, Familia, Juzgado, etc.)

	// 3. Historial de Consumo
	EdadInicioConsumo int         `json:"edad_inicio"`
	SustanciaPrincipal string      `json:"sustancia_principal"`
	FrecuenciaConsumo  string      `json:"frecuencia"`         // Diario, Semanal, Ocasional
	ViaAdministracion  string      `json:"via_admin"`          // Oral, Nasal, Fumada, Inyectada
	Sustancias         []Sustancia `gorm:"many2many:perfil_sustancias;" json:"sustancias"`

	// 4. Salud y Tratamiento
	AntecedentesMedicos      string `json:"ant_medicos"`
	AntecedentesPsicologicos string `json:"ant_psicologicos"`
	Alergias                 string `json:"alergias"`
	MedicacionActual         string `json:"medicacion"`
	TratamientosAnteriores   string `json:"trat_anteriores"`

	// 5. Situación Social y Legal
	SituacionVivienda      string `json:"vivienda"`            // Propia, Alquiler, Familia, Sin techo
	ProblemasLegales       string `json:"prob_legales"`
	ContactoEmergenciaNom  string `json:"emergencia_nom"`
	ContactoEmergenciaTel  string `json:"emergencia_tel"`
	ParentescoEmergencia   string `json:"emergencia_parentesco"`

	Observaciones     string      `json:"observaciones"`
}

func (p *Perfil) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
