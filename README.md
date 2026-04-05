# Gestión de Drogodependencia - Andalucía

> **Reemplazo programado con IA y revisado por mí, el sistema de drogodependencia de la junta es un poco meh, esto intenta hacer ver que con 1h y inteligencia artificial hay un reemplazo open-source.**

Este proyecto es una aplicación web moderna y eficiente desarrollada en Go para la gestión integral de expedientes clínicos en centros de tratamiento de drogodependencia, diseñada con una estética profesional y funcional.

## Características Principales

- **Gestión de Pacientes:** Búsqueda rápida por DNI o Pasaporte para generar nuevos informes sobre pacientes existentes.
- **Registro Completo:** Formulario detallado para nuevos ingresos que incluye:
  - Datos personales y situación socio-laboral.
  - Historial de ingreso y derivación.
  - Perfil toxicológico avanzado con buscador de sustancias en tiempo real.
  - Evaluación de salud física y mental.
  - Entorno social y legal.
- **Validación Estricta:** Control de formatos para fechas (DD/MM/AAAA), DNI/Pasaporte y correos electrónicos.
- **Generación de Informes:** Exportación instantánea a PDF profesional con el branding de la Junta de Andalucía.
- **Arquitectura Moderna:**
  - **Backend:** Go + Fiber + GORM (SQLite).
  - **Frontend:** HTMX + Templ + Tailwind CSS (Interactivo y sin dependencias pesadas de JS).
  - **Reportes:** Maroto v2 para PDFs de alta calidad.

---

## 🚀 Mini-Tutorial de Inicio Rápido

### 1. Requisitos Previos
Asegúrate de tener instalado:
- **Go 1.21+**: [Instalar Go](https://go.dev/doc/install)
- **Make**: Para ejecutar los comandos de automatización.

### 2. Instalación y Configuración
Clona el repositorio y entra en la carpeta:
```bash
git clone https://github.com/lilsheepyy/drogodependencia-andalucia.git
cd drogodependencia-andalucia
```

Instala la herramienta de plantillas `templ`:
```bash
go install github.com/a-h/templ/cmd/templ@latest
```

### 3. Generación de Código y Construcción
Genera las plantillas HTML y compila el binario:
```bash
make build
```

### 4. Ejecución
Inicia el servidor local:
```bash
make run
```
La aplicación estará disponible en: [**http://localhost:8080**](http://localhost:8080)

### 5. Comandos Útiles
- `make generate`: Solo regenera las plantillas `templ`.
- `make clean`: Borra los binarios, la base de datos local y archivos temporales.

---

## Estructura del Proyecto

- `cmd/server/`: Punto de entrada del servidor Fiber.
- `internal/handlers/`: Controladores de las rutas y lógica HTMX.
- `internal/models/`: Definiciones de base de datos (GORM).
- `internal/repository/`: Capa de persistencia (SQLite).
- `internal/service/`: Lógica de negocio (Generación de PDF).
- `internal/views/`: Componentes de UI declarativos con Templ.
- `static/`: Activos estáticos (Logo, CSS personalizado).

## Licencia
Open Source - Siente libre de mejorar y adaptar este sistema.
