package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/services" // 🌟 Corregido: Se añade la importación de servicios

	"github.com/go-chi/chi/v5"
)

type MecanicosHandler struct {
	service *services.MecanicosService // 🌟 Corregido: Nombre semántico correcto 'service'
}

// NewMecanicosHandler ahora recibe estrictamente el Servicio inyectado desde el main.go
func NewMecanicosHandler(s *services.MecanicosService) *MecanicosHandler {
	return &MecanicosHandler{service: s}
}

// Routes define los endpoints para el módulo de Mecánicos
func (h *MecanicosHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListarMecanicos)     // GET /api/v1/mecanicos
	r.Post("/", h.RegistrarMecanico)  // POST /api/v1/mecanicos
	r.Get("/{id}", h.ObtenerMecanico) // GET /api/v1/mecanicos/{id}

	return r
}

// GET /api/v1/mecanicos
func (h *MecanicosHandler) ListarMecanicos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 🌟 Corregido: Llama al método del SERVICIO (revisar que en su service se llame ListarMecanicos)
	mecanicos, err := h.service.ListarMecanicos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error al obtener la lista de mecánicos"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mecanicos)
}

// POST /api/v1/mecanicos
func (h *MecanicosHandler) RegistrarMecanico(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nuevoMecanico models.Mecanico
	if err := json.NewDecoder(r.Body).Decode(&nuevoMecanico); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "El cuerpo de la petición JSON no es válido"}`))
		return
	}

	// 🌟 Corregido: Pasa la estructura al SERVICIO para procesar las reglas de negocio
	resultado, err := h.service.RegistrarMecanico(nuevoMecanico)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resultado)
}

// GET /api/v1/mecanicos/{id}
func (h *MecanicosHandler) ObtenerMecanico(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr) // 🌟 Aquí id ya es de tipo 'int'
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "El ID del mecánico debe ser un número entero válido"}`))
		return
	}

	// 🌟 Le pasamos el int directamente sin rodeos
	mecanico, err := h.service.ObtenerMecanicoPorID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "El mecánico solicitado no existe en la base de datos"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mecanico)
}
