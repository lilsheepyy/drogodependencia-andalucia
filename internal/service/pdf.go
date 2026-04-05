package service

import (
	"fmt"
	"time"

	"github.com/drugprofile/drugprofile/internal/models"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

type PDFService struct{}

func NewPDFService() *PDFService {
	return &PDFService{}
}

func (s *PDFService) GeneratePerfilPDF(p *models.Perfil) ([]byte, error) {
	cfg := config.NewBuilder().
		WithPageNumber().
		WithLeftMargin(15).
		WithRightMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(cfg)

	// Colores y Estilos
	blueColor := &props.Color{Red: 30, Green: 58, Blue: 138}
	grayColor := &props.Color{Red: 71, Green: 85, Blue: 105}
	lightGrayColor := &props.Color{Red: 226, Green: 232, Blue: 240}

	sectionTitleProp := props.Text{Size: 10, Style: fontstyle.Bold, Color: blueColor}
	labelProp := props.Text{Size: 7, Style: fontstyle.Bold, Color: grayColor}
	valueProp := props.Text{Size: 8, Style: fontstyle.Normal, Top: 3}

	// 1. CABECERA
	m.AddRows(
		row.New(15).Add(
			col.New(2).Add(
				image.NewFromFile("static/logo.svg", props.Rect{
					Center: true,
					Percent: 100,
				}),
			),
			col.New(6).Add(text.New("GESTIÓN DE DROGODEPENDENCIA", props.Text{Size: 14, Style: fontstyle.Bold, Color: blueColor, Top: 2})),
			col.New(4).Add(text.New("EXPEDIENTE CLÍNICO", props.Text{Size: 9, Style: fontstyle.Bold, Align: align.Right})),
		),
		row.New(6).Add(
			col.New(2),
			col.New(6).Add(text.New("Gestión Hospitalaria de Drogodependencias", props.Text{Size: 7, Style: fontstyle.Italic, Color: grayColor})),
			col.New(4).Add(text.New(fmt.Sprintf("Expediente: %s", p.ID[:8]), props.Text{Size: 7, Align: align.Right})),
		),
	)

	m.AddRows(row.New(4).Add(col.New(12).Add(line.New(props.Line{Color: blueColor, Thickness: 0.5}))))
	m.AddRows(row.New(4))

	// 2. DATOS PERSONALES
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("I. INFORMACIÓN DEL PACIENTE", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))
	
	m.AddRows(
		row.New(11).Add(
			col.New(6).Add(text.New("Nombre Completo", labelProp), text.New(p.NombreCompleto, valueProp)),
			col.New(3).Add(text.New("DNI/Pasaporte", labelProp), text.New(p.DNI, valueProp)),
			col.New(3).Add(text.New("F. Nacimiento", labelProp), text.New(p.FechaNacimiento.Format("02/01/2006"), valueProp)),
		),
		row.New(11).Add(
			col.New(6).Add(text.New("Dirección", labelProp), text.New(p.Direccion, valueProp)),
			col.New(3).Add(text.New("Teléfono", labelProp), text.New(p.Telefono, valueProp)),
			col.New(3).Add(text.New("Email", labelProp), text.New(p.Email, valueProp)),
		),
		row.New(11).Add(
			col.New(4).Add(text.New("Estado Civil", labelProp), text.New(p.EstadoCivil, valueProp)),
			col.New(4).Add(text.New("Situación Laboral", labelProp), text.New(p.SituacionLaboral, valueProp)),
			col.New(4).Add(text.New("Nivel de Estudios", labelProp), text.New(p.NivelEstudios, valueProp)),
		),
	)

	// 3. ADMISIÓN E INGRESO
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("II. DETALLES DE ADMISIÓN", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))
	
	m.AddRows(
		row.New(11).Add(
			col.New(4).Add(text.New("Fecha Ingreso", labelProp), text.New(p.FechaIngreso.Format("02/01/2006"), valueProp)),
			col.New(8).Add(text.New("Derivado por / Referencia", labelProp), text.New(p.DerivadoPor, valueProp)),
		),
		row.New(12).Add(
			col.New(12).Add(text.New("Motivo Principal", labelProp), text.New(p.MotivoIngreso, valueProp)),
		),
	)

	// 4. HISTORIAL TOXICOLÓGICO
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("III. PERFIL TOXICOLÓGICO", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))

	m.AddRows(
		row.New(11).Add(
			col.New(3).Add(text.New("Edad Inicio", labelProp), text.New(fmt.Sprintf("%d años", p.EdadInicioConsumo), valueProp)),
			col.New(3).Add(text.New("Sustancia Principal", labelProp), text.New(p.SustanciaPrincipal, valueProp)),
			col.New(3).Add(text.New("Frecuencia", labelProp), text.New(p.FrecuenciaConsumo, valueProp)),
			col.New(3).Add(text.New("Vía Administración", labelProp), text.New(p.ViaAdministracion, valueProp)),
		),
	)

	if len(p.Sustancias) > 0 {
		sustanciaText := ""
		for i, s := range p.Sustancias {
			if i > 0 { sustanciaText += ", " }
			sustanciaText += s.Nombre
		}
		m.AddRows(row.New(11).Add(col.New(12).Add(text.New("Otras sustancias identificadas", labelProp), text.New(sustanciaText, valueProp))))
	}

	// 5. EVALUACIÓN CLÍNICA
	m.AddRows(row.New(5).Add(col.New(12).Add(text.New("IV. EVALUACIÓN CLÍNICA Y SALUD", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))

	m.AddRows(
		row.New(11).Add(col.New(6).Add(text.New("Antecedentes Médicos", labelProp), text.New(p.AntecedentesMedicos, valueProp)),
		                col.New(6).Add(text.New("Salud Mental", labelProp), text.New(p.AntecedentesPsicologicos, valueProp))),
		row.New(10).Add(col.New(6).Add(text.New("Alergias Conocidas", labelProp), text.New(p.Alergias, valueProp)),
		                col.New(6).Add(text.New("Medicación Actual", labelProp), text.New(p.MedicacionActual, valueProp))),
		row.New(10).Add(col.New(12).Add(text.New("Tratamientos de Rehabilitación Previos", labelProp), text.New(p.TratamientosAnteriores, valueProp))),
	)

	// 6. ENTORNO SOCIAL Y LEGAL
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("V. SITUACIÓN SOCIAL Y LEGAL", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))

	m.AddRows(
		row.New(11).Add(
			col.New(4).Add(text.New("Vivienda", labelProp), text.New(p.SituacionVivienda, valueProp)),
			col.New(8).Add(text.New("Asuntos Legales", labelProp), text.New(p.ProblemasLegales, valueProp)),
		),
		row.New(11).Add(
			col.New(12).Add(
				text.New("Contacto de Emergencia", labelProp),
				text.New(fmt.Sprintf("%s (%s) — Tel: %s", p.ContactoEmergenciaNom, p.ParentescoEmergencia, p.ContactoEmergenciaTel), valueProp),
			),
		),
	)

	// 7. OBSERVACIONES (Con límite de seguridad para página única)
	m.AddRows(row.New(6).Add(col.New(12).Add(text.New("VI. OBSERVACIONES FINALES", sectionTitleProp))))
	m.AddRows(row.New(1).Add(col.New(12).Add(line.New(props.Line{Color: lightGrayColor, Thickness: 0.2}))))
	
	obs := p.Observaciones
	if len(obs) > 350 {
		obs = obs[:347] + "..."
	}
	m.AddRows(row.New(15).Add(col.New(12).Add(text.New(obs, valueProp))))

	// PIE DE PÁGINA
	m.AddRows(
		row.New(15).Add(
			col.New(12).Add(
				line.New(props.Line{Color: grayColor, Thickness: 0.1}),
				text.New(fmt.Sprintf("Generado el %s — Copia Certificada — Confidencial", time.Now().Format("02/01/2006 15:04")), props.Text{
					Top:   6,
					Size:  6,
					Align: align.Center,
					Color: grayColor,
				}),
			),
		),
	)

	document, err := m.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), nil
}
