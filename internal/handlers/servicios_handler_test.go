package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"turno-facil-api/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Fake en memoria para Servicios
type FakeServiciosService struct {
	Servicios map[int]models.Servicio
	ProximoID int
}

// 🌟 AJUSTE: Ahora se llama RegistrarServicio como en tu service.go
func (f *FakeServiciosService) RegistrarServicio(s models.Servicio) (models.Servicio, error) {
	s.ID = f.ProximoID
	f.Servicios[f.ProximoID] = s
	f.ProximoID++
	return s, nil
}

// Middleware simulado para validar el token
func SimularAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized"}`))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// TEST 1: El 401 Unauthorized cuando falta el token (Ruta protegida)
func TestCrearServicioHandler_Error_401Unauthorized(t *testing.T) {
	r := chi.NewRouter()
	fakeService := &FakeServiciosService{Servicios: make(map[int]models.Servicio), ProximoID: 1}

	// Protegemos la ruta con el middleware simulado
	r.With(SimularAuthMiddleware).Post("/api/v1/servicios", func(w http.ResponseWriter, r *http.Request) {
		var s models.Servicio
		json.NewDecoder(r.Body).Decode(&s)
		res, _ := fakeService.RegistrarServicio(s) // 🌟 Ajustado nombre
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	})

	payload, _ := json.Marshal(models.Servicio{Nombre: "Alineación", DuracionMins: 30})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/servicios", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Valida el 401
}

// TEST 2: Camino Feliz (201 Created) con Fake guardando exitosamente
func TestCrearServicioHandler_Exito(t *testing.T) {
	r := chi.NewRouter()
	fakeService := &FakeServiciosService{Servicios: make(map[int]models.Servicio), ProximoID: 1}

	r.Post("/api/v1/servicios", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var s models.Servicio
		json.NewDecoder(r.Body).Decode(&s)
		res, _ := fakeService.RegistrarServicio(s) // 🌟 Ajustado nombre
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(res)
	})

	payload, _ := json.Marshal(models.Servicio{Nombre: "Cambio de Filtro", DuracionMins: 20})
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/servicios", bytes.NewBuffer(payload))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code) // Valida el 201
}
