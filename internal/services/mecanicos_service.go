package services

import (
	"errors"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
)

type MecanicosService struct {
	repo repository.MecanicosRepository
}

func NewMecanicosService(repo repository.MecanicosRepository) *MecanicosService {
	return &MecanicosService{
		repo: repo,
	}
}

func (s *MecanicosService) ListarMecanicos() ([]models.Mecanico, error) {
	return s.repo.Listar()
}

// 🌟 Todo con int estándar
func (s *MecanicosService) ObtenerMecanicoPorID(id int) (models.Mecanico, error) {
	if id <= 0 {
		return models.Mecanico{}, errors.New("ID inválido")
	}
	return s.repo.ObtenerPorID(id)
}

func (s *MecanicosService) RegistrarMecanico(m models.Mecanico) (models.Mecanico, error) {
	if m.Nombre == "" {
		return models.Mecanico{}, errors.New("el nombre del mecánico es obligatorio")
	}

	err := s.repo.Crear(&m)
	if err != nil {
		return models.Mecanico{}, err
	}

	return m, nil
}

// 🌟 CambiarDisponibilidad llama internamente a ActualizarEstado usando int
func (s *MecanicosService) CambiarDisponibilidad(id int, activo bool) error {
	return s.repo.CambiarDisponibilidad(id, activo)
}
