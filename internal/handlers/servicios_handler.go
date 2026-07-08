package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/services"

	"github.com/go-chi/chi/v5"
)

type ServiciosHandler struct {
	service *services.ServiciosService // 🌟 Cambiado de repo a service
}

// NewServiciosHandler recibe el servicio inyectado desde el main.go
func NewServiciosHandler(s *services.ServiciosService) *ServiciosHandler {
	return &ServiciosHandler{service: s}
}

// Routes define los endpoints para el módulo de Servicios
func (h *ServiciosHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.ListarServicios)     // GET /api/v1/servicios
	r.Post("/", h.CrearServicio)      // POST /api/v1/servicios
	r.Get("/{id}", h.ObtenerServicio) // GET /api/v1/servicios/{id}

	return r
}

// GET /api/v1/servicios
func (h *ServiciosHandler) ListarServicios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 🌟 CORRECCIÓN: Llama al método del SERVICIO
	servicios, err := h.service.ListarServicios()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error al obtener la lista de servicios"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servicios)
}

// POST /api/v1/servicios
func (h *ServiciosHandler) CrearServicio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nuevoServicio models.Servicio
	if err := json.NewDecoder(r.Body).Decode(&nuevoServicio); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "JSON inválido o mal estructurado"}`))
		return
	}

	// 🌟 CORRECCIÓN: Ahora delegamos la creación y sus validaciones al SERVICIO
	servicioCreado, err := h.service.RegistrarServicio(nuevoServicio)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		// Le mostramos al cliente el error de negocio dinámico que devuelva el servicio
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(servicioCreado)
}

// GET /api/v1/servicios/{id}
func (h *ServiciosHandler) ObtenerServicio(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "El ID proporcionado debe ser un número entero válido"}`))
		return
	}

	// 🌟 CORRECCIÓN: Si tu servicio tiene el método para buscar por ID (puedes añadirlo si falta),
	// se llama a través de h.service. De lo contrario, asegúrate de tenerlo en servicios_service.go
	servicio, err := h.service.ObtenerServicioPorID(int(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Servicio técnico no encontrado en la base de datos"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(servicio)
}
