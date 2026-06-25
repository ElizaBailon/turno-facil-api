package handlers

import (
	"context"
	"net/http"
	"strings"
	"turno-facil-api/internal/services"

	"github.com/golang-jwt/jwt/v5"
)

// Definimos un tipo único para las llaves del contexto y evitar colisiones
type contextKey string

const RolKey contextKey = "rol"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Se requiere token de autenticación (Header Authorization vacío)", http.StatusUnauthorized)
			return
		}

		// El header suele venir como "Bearer <token>"
		partes := strings.Split(authHeader, " ")
		if len(partes) != 2 || partes[0] != "Bearer" {
			http.Error(w, "Formato de token inválido. Debe ser 'Bearer <token>'", http.StatusUnauthorized)
			return
		}

		tokenString := partes[1]

		// Validar el token JWT
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validar el método de firma
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return services.ClaveSecreta, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido o expirado", http.StatusUnauthorized)
			return
		}

		// Si el token es válido, extraemos los claims (datos del usuario)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			rol := claims["rol"].(string)
			// Guardamos el rol en el contexto por si una ruta necesita validar si es "admin"
			ctx := context.WithValue(r.Context(), RolKey, rol)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRol verifica si el usuario tiene el rol necesario para acceder al endpoint
func RequireRol(rolPermitido string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraemos el rol del contexto (guardado previamente por AuthMiddleware)
			rolUsuario, ok := r.Context().Value(RolKey).(string)

			// Si no hay rol o no coincide con el requerido, bloqueamos con 403 Forbidden
			if !ok || rolUsuario != rolPermitido {
				http.Error(w, "Acceso denegado: no tienes los permisos requeridos para esta acción", http.StatusForbidden)
				return
			}

			// Si el rol es correcto, continúa al siguiente Handler
			next.ServeHTTP(w, r)
		})
	}
}
