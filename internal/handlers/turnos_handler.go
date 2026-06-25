package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/services"

	"github.com/go-chi/chi/v5"
)

type TurnosHandler struct {
	service *services.TurnosService
}

func NewTurnosHandler(s *services.TurnosService) *TurnosHandler {
	return &TurnosHandler{service: s}
}

// Routes define los endpoints para el módulo de Turnos
func (h *TurnosHandler) Routes() chi.Router {
	r := chi.NewRouter()

	// 🌟 Cualquier usuario autenticado puede ver la lista de turnos
	r.Get("/", h.ListarTurnos)

	// 🌟 PROTECCIÓN DE ROL: Solo el rol "admin" puede crear turnos
	r.With(RequireRol("admin")).Post("/", h.CrearTurno)

	// Endpoint para eliminar turnos por ID
	r.Delete("/{id}", h.EliminarTurno)

	r.With(RequireRol("admin")).Put("/{id}", h.ActualizarTurno)

	return r
}

// GET /api/v1/turnos
func (h *TurnosHandler) ListarTurnos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	turnos, err := h.service.ObtenerTodos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error al listar turnos"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(turnos)
}

// POST /api/v1/turnos
func (h *TurnosHandler) CrearTurno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var t models.Turno
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "JSON inválido o mal estructurado"}`))
		return
	}

	if t.VehiculoID <= 0 || t.MecanicoID <= 0 || t.ServicioID <= 0 || t.FechaHora.IsZero() {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Campos obligatorios incompletos"}`))
		return
	}

	t.Estado = "pendiente"
	t.DuracionEst = 60 // Duración por defecto o heredada del servicio

	turnoGuardado, err := h.service.RegistrarTurno(t)
	if err != nil {
		w.WriteHeader(http.StatusConflict) // 409 Conflict para choques de horarios
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(turnoGuardado)
}

// DELETE /api/v1/turnos/{id}
func (h *TurnosHandler) EliminarTurno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "ID de turno inválido"}`))
		return
	}

	if err := h.service.Eliminar(id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "No se pudo eliminar el turno"}`))
		return
	}

	w.WriteHeader(http.StatusOK) // 200
	w.Write([]byte(`{"message": "Turno eliminado exitosamente"}`))
}

// PUT /api/v1/turnos/{id}
func (h *TurnosHandler) ActualizarTurno(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Obtener el ID de la URL
	idStr := chi.URLParam(r, "id")
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de turno inválido"})
		return
	}

	// 2. Decodificar el cuerpo (solo los campos editables)
	var t models.Turno
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "JSON mal estructurado"})
		return
	}

	// 3. Llamar al servicio
	turnoActualizado, err := h.service.Actualizar(id, t)
	if err != nil {
		// Si el error es por conflicto de horario, enviamos 409
		if err.Error() == "conflicto de agenda: el mecánico seleccionado ya está ocupado en ese horario" {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	// 4. Responder con el turno completo (rellenado con Preloads)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(turnoActualizado)
}
