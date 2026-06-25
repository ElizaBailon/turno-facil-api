package services

import (
	"errors"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
)

type ServiciosService struct {
	repo repository.ServiciosRepository
}

// NewServiciosService inyecta el repositorio correspondiente
func NewServiciosService(repo repository.ServiciosRepository) *ServiciosService {
	return &ServiciosService{repo: repo}
}

// ListarServicios obtiene el catálogo completo del taller
func (s *ServiciosService) ListarServicios() ([]models.Servicio, error) {
	return s.repo.Listar()
}

// RegistrarServicio añade un nuevo trabajo al taller con validaciones de negocio
func (s *ServiciosService) RegistrarServicio(serv models.Servicio) (models.Servicio, error) {
	if serv.Nombre == "" {
		return models.Servicio{}, errors.New("el nombre del servicio no puede estar vacío")
	}
	if serv.DuracionMins <= 0 {
		return models.Servicio{}, errors.New("la duración estimada debe ser mayor a 0 minutos")
	}

	err := s.repo.Crear(&serv)
	return serv, err
}

func (s *ServiciosService) ObtenerServicioPorID(id int) (models.Servicio, error) {
	return s.repo.ObtenerPorID(id)
}
