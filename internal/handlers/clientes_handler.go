package handlers

import (
	"encoding/json"
	"net/http"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/services"

	"github.com/go-chi/chi/v5"
)

type ClientesHandler struct {
	service *services.ClientesService
}

func NewClientesHandler(s *services.ClientesService) *ClientesHandler {
	return &ClientesHandler{service: s}
}

func (h *ClientesHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.RegistrarCliente)
	r.Get("/{cedula}", h.BuscarCliente)
	r.Put("/{cedula}", h.ActualizarCliente)
	r.Delete("/{id}", h.EliminarCliente)
	return r
}

func (h *ClientesHandler) RegistrarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var c models.Cliente
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "JSON inválido"}`))
		return
	}

	if c.Cedula == "" || c.Nombre == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Cédula y nombre son obligatorios"}`))
		return
	}

	if err := h.service.RegistrarCliente(&c); err != nil {
		if err.Error() == "duplicado" {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error": "El cliente con esta cédula ya se encuentra registrado"}`))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "No se pudo guardar el cliente"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *ClientesHandler) BuscarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cedula := chi.URLParam(r, "cedula")

	cliente, err := h.service.ObtenerPorCedula(cedula)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Cliente no encontrado"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cliente)
}

func (h *ClientesHandler) ActualizarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cedula := chi.URLParam(r, "cedula")

	cliente, err := h.service.ObtenerPorCedula(cedula)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Cliente no encontrado"}`))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cliente); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "JSON inválido"}`))
		return
	}

	cliente.Cedula = cedula

	if err := h.service.Actualizar(&cliente); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "No se pudo actualizar el cliente"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cliente)
}

func (h *ClientesHandler) EliminarCliente(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	if err := h.service.Eliminar(id); err != nil {
		// 🌟 Si el error dice explícitamente "cliente no encontrado" devolvemos un 404
		if err.Error() == "cliente no encontrado" || err.Error() == "ID de cliente inválido" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		}

		// 🌟 Si falla por restricciones relacionales (llaves foráneas), responde un 500 descriptivo
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "No se pudo eliminar el cliente. Verifique que no tenga vehículos o turnos asociados."}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Cliente eliminado exitosamente"}`))
}
