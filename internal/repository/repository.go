package repository

import (
	"encoding/json"
	"os"

	"github.com/drugprofile/drugprofile/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dbPath string) (*Repository, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Sustancia{}, &models.Perfil{}, &models.User{})
	if err != nil {
		return nil, err
	}

	repo := &Repository{db: db}
	repo.seedSustancias()
	repo.seedUsers()

	return repo, nil
}

func (r *Repository) seedUsers() {
	var count int64
	r.db.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	// Crear admin:admin
	hashedAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	r.db.Create(&models.User{
		DNI:                 "admin",
		Password:            string(hashedAdmin),
		IsAdmin:             true,
		NeedsPasswordChange: true,
	})

	// Crear user:user
	hashedUser, _ := bcrypt.GenerateFromPassword([]byte("user"), bcrypt.DefaultCost)
	r.db.Create(&models.User{
		DNI:                 "user",
		Password:            string(hashedUser),
		IsAdmin:             false,
		NeedsPasswordChange: true,
	})
}

func (r *Repository) GetUserByDNI(dni string) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, "dni = ?", dni).Error
	return &u, err
}

func (r *Repository) GetUserByID(id string) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, "id = ?", id).Error
	return &u, err
}

func (r *Repository) CreateUser(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) DeleteUser(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

func (r *Repository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Order("created_at desc").Find(&users).Error
	return users, err
}

func (r *Repository) UpdateUser(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *Repository) seedSustancias() {
	file, err := os.ReadFile("drogas.json")
	if err != nil {
		return
	}

	var sustancias []string
	if err := json.Unmarshal(file, &sustancias); err != nil {
		return
	}

	for _, s := range sustancias {
		var count int64
		// Solo contamos las que NO son manuales
		r.db.Model(&models.Sustancia{}).Where("nombre = ? AND es_manual = ?", s, false).Count(&count)
		if count == 0 {
			r.db.Create(&models.Sustancia{Nombre: s, EsManual: false})
		}
	}
}

// GetSustancias solo retorna las sustancias fijas (no manuales)
func (r *Repository) GetSustancias(search string) ([]models.Sustancia, error) {
	var sustancias []models.Sustancia
	query := r.db.Model(&models.Sustancia{}).Where("es_manual = ?", false)
	if search != "" {
		query = query.Where("nombre LIKE ?", "%"+search+"%")
	}
	err := query.Find(&sustancias).Error
	return sustancias, err
}

func (r *Repository) CreateSustanciaManual(nombre string) (*models.Sustancia, error) {
	// Intentamos ver si ya existe como manual o fija
	var s models.Sustancia
	err := r.db.Where("nombre = ?", nombre).First(&s).Error
	if err == nil {
		// Ya existe, la retornamos tal cual (aunque sea fija, se puede usar como manual)
		return &s, nil
	}

	// No existe, la creamos marcada como manual
	s = models.Sustancia{Nombre: nombre, EsManual: true}
	err = r.db.Create(&s).Error
	return &s, err
}

func (r *Repository) CreateSustancia(nombre string) (*models.Sustancia, error) {
	s := &models.Sustancia{Nombre: nombre, EsManual: false}
	err := r.db.Create(s).Error
	return s, err
}

func (r *Repository) CreatePerfil(p *models.Perfil) error {
	return r.db.Create(p).Error
}

func (r *Repository) GetPerfil(id string) (*models.Perfil, error) {
	var p models.Perfil
	err := r.db.Preload("Sustancias").First(&p, "id = ?", id).Error
	return &p, err
}

func (r *Repository) GetPerfilByDNI(dni string) (*models.Perfil, error) {
	var p models.Perfil
	err := r.db.Preload("Sustancias").Order("created_at DESC").First(&p, "dni = ?", dni).Error
	return &p, err
}

func (r *Repository) GetSustanciaByNombre(nombre string) (*models.Sustancia, error) {
	var s models.Sustancia
	err := r.db.Where("nombre = ?", nombre).First(&s).Error
	return &s, err
}
