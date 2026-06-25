package services

import (
	"errors"
	"time"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"
)

// TurnosServiceI define el contrato limpio de la interfaz
type TurnosServiceI interface {
	RegistrarTurno(nuevo models.Turno) (models.Turno, error)
	ObtenerTodos() ([]models.Turno, error)
	Eliminar(id int) error // 🌟 Corregido a 'int' para mantener consistencia
	Actualizar(id int, turnoActualizado models.Turno) (models.Turno, error)
	ObtenerPorID(id int) (models.Turno, error)
}

type TurnosService struct {
	repo         repository.TurnosRepository
	servicioRepo repository.ServiciosRepository // 🌟 Inyectamos el repositorio de servicios
}

func NewTurnosService(repo repository.TurnosRepository, servicioRepo repository.ServiciosRepository) *TurnosService {
	return &TurnosService{
		repo:         repo,
		servicioRepo: servicioRepo,
	}
}

func (s *TurnosService) RegistrarTurno(nuevo models.Turno) (models.Turno, error) {
	// 🌟 PASO OBLIGATORIO: Buscar el servicio en la base de datos para heredar su duración real
	servicioReal, err := s.servicioRepo.ObtenerPorID(nuevo.ServicioID)
	if err != nil {
		return models.Turno{}, errors.New("el servicio seleccionado no existe en el catálogo")
	}

	// Asignamos la verdadera duración configurada en el modelo de servicios
	nuevo.DuracionEst = servicioReal.DuracionMins

	if nuevo.DuracionEst <= 0 {
		return models.Turno{}, errors.New("la duración estimada del servicio debe ser mayor a 0 minutos")
	}

	// Consultar los turnos que ya tiene agendados ese mecánico específico
	turnosExistentes, err := s.repo.ListarPorMecanico(nuevo.MecanicoID)
	if err != nil {
		return models.Turno{}, err
	}

	inicioNuevo := nuevo.FechaHora
	finNuevo := inicioNuevo.Add(time.Duration(nuevo.DuracionEst) * time.Minute)

	// Algoritmo de detección de choques de horarios
	for _, t := range turnosExistentes {
		if t.Estado == "cancelado" || t.Estado == "listo" {
			continue
		}
		inicioExistente := t.FechaHora
		finExistente := inicioExistente.Add(time.Duration(t.DuracionEst) * time.Minute)

		if inicioNuevo.Before(finExistente) && finNuevo.After(inicioExistente) {
			return models.Turno{}, errors.New("conflicto de agenda: el mecánico seleccionado ya está ocupado en ese horario")
		}
	}

	err = s.repo.Crear(&nuevo)
	return nuevo, err
}

func (s *TurnosService) ObtenerTodos() ([]models.Turno, error) {
	return s.repo.Listar()
}

func (s *TurnosService) Eliminar(id int) error {
	if id <= 0 {
		return errors.New("ID de turno inválido")
	}
	return s.repo.Eliminar(id)
}

func (s *TurnosService) Actualizar(id int, turnoActualizado models.Turno) (models.Turno, error) {
	// 1. Verificar que el turno original existe
	turno, err := s.repo.BuscarPorID(id)
	if err != nil {
		return models.Turno{}, errors.New("el turno no existe")
	}

	// 2. Actualizar los campos
	turno.FechaHora = turnoActualizado.FechaHora
	turno.MecanicoID = turnoActualizado.MecanicoID
	turno.ServicioID = turnoActualizado.ServicioID
	turno.Notas = turnoActualizado.Notas

	// Recalcular duración
	servicioReal, _ := s.servicioRepo.ObtenerPorID(turno.ServicioID)
	turno.DuracionEst = servicioReal.DuracionMins

	// 3. 🌟 VALIDACIÓN REUTILIZADA:
	// Pasamos el ID del turno actual (turno.ID) para que la función lo IGNORE
	// durante la comparación de horarios.
	if err := s.esHorarioDisponible(turno.MecanicoID, turno.FechaHora, turno.DuracionEst, turno.ID); err != nil {
		return models.Turno{}, err
	}

	// 4. Guardar cambios
	err = s.repo.Actualizar(&turno)

	return s.repo.BuscarPorID(id)
}

// 1. Nueva función privada de validación
func (s *TurnosService) esHorarioDisponible(mecanicoID int, inicio time.Time, duracion int, turnoIDIgnorar int) error {
	turnosExistentes, err := s.repo.ListarPorMecanico(mecanicoID)
	if err != nil {
		return err
	}

	finNuevo := inicio.Add(time.Duration(duracion) * time.Minute)

	for _, t := range turnosExistentes {
		// Ignorar el turno actual si estamos editando (usamos el ID)
		if t.ID == turnoIDIgnorar || t.Estado == "cancelado" || t.Estado == "listo" {
			continue
		}

		inicioExistente := t.FechaHora
		finExistente := inicioExistente.Add(time.Duration(t.DuracionEst) * time.Minute)

		if inicio.Before(finExistente) && finNuevo.After(inicioExistente) {
			return errors.New("conflicto de agenda: el mecánico seleccionado ya está ocupado en ese horario")
		}
	}
	return nil
}

func (s *TurnosService) ObtenerPorID(id int) (models.Turno, error) {
	if id <= 0 {
		return models.Turno{}, errors.New("ID de turno inválido")
	}
	return s.repo.BuscarPorID(id)
}
