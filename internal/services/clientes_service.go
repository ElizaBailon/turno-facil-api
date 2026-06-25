package services

import (
	"errors"
	"strconv"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
)

type ClientesService struct {
	repo repository.ClientesRepository
}

func NewClientesService(repo repository.ClientesRepository) *ClientesService {
	return &ClientesService{repo: repo}
}

func (s *ClientesService) RegistrarCliente(c *models.Cliente) error {
	if c.Nombre == "" || c.Cedula == "" {
		return errors.New("el nombre y la cédula son campos obligatorios")
	}
	if len(c.Cedula) != 10 {
		return errors.New("la cédula debe tener exactamente 10 dígitos")
	}
	return s.repo.Crear(c)
}

func (s *ClientesService) ObtenerPorCedula(cedula string) (models.Cliente, error) {
	if cedula == "" {
		return models.Cliente{}, errors.New("la cédula es requerida")
	}
	return s.repo.BuscarPorCedula(cedula)
}

// Actualizar reenvía los datos modificados del cliente hacia el repositorio
func (s *ClientesService) Actualizar(c *models.Cliente) error {
	if c.Nombre == "" || c.Cedula == "" {
		return errors.New("el nombre y la cédula no pueden estar vacíos")
	}
	return s.repo.Actualizar(c)
}

func (s *ClientesService) Eliminar(idStr string) error {
	// Convertimos el ID de texto que viene de Chi a un entero estándar
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return errors.New("ID de cliente inválido")
	}

	// Delegamos la eliminación física a la capa de persistencia
	return s.repo.Eliminar(id)
}
