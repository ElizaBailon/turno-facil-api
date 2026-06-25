package services

import (
	"errors"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
)

// MecanicosService manejará la lógica de negocio del Integrante A
type MecanicosService struct {
	repo repository.MecanicosRepository
}

// NewMecanicosService crea una nueva instancia inyectando el repositorio requerido
func NewMecanicosService(repo repository.MecanicosRepository) *MecanicosService {
	return &MecanicosService{
		repo: repo,
	}
}

func (s *MecanicosService) ListarMecanicos() ([]models.Mecanico, error) {
	return s.repo.Listar()
}

func (s *MecanicosService) ObtenerMecanicoPorID(id int) (models.Mecanico, error) {
	if id <= 0 {
		return models.Mecanico{}, errors.New("ID inválido")
	}
	return s.repo.ObtenerPorID(id) // Reenvía el int al repositorio
}

func (s *MecanicosService) RegistrarMecanico(m models.Mecanico) (models.Mecanico, error) {
	if m.Nombre == "" {
		return models.Mecanico{}, errors.New("el nombre del mecánico es obligatorio")
	}

	// 🌟 Conectado al repo real: Guardamos en la base de datos
	err := s.repo.Crear(&m)
	if err != nil {
		return models.Mecanico{}, err
	}

	return m, nil
}

func (s *MecanicosService) CambiarDisponibilidad(id int, activo bool) error {
	// 🌟 Conectado al repo real: Ahora sí altera el estado en SQLite de verdad
	return s.repo.CambiarDisponibilidad(id, activo)
}
