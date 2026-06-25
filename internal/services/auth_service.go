package services

import (
	"errors"
	"time"
	"turno-facil-api/internal/models"
	"turno-facil-api/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ClaveSecreta = []byte("mi_clave_secreta_super_segura_para_el_taller")

type AuthService struct {
	repo repository.UsuariosRepository
}

func NewAuthService(repo repository.UsuariosRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Registrar guarda al usuario con la contraseña hasheada
func (s *AuthService) Registrar(username, password, rol string) (models.Usuario, error) {
	if username == "" || password == "" {
		return models.Usuario{}, errors.New("el usuario y la contraseña son obligatorios")
	}

	// 🔐 Generamos el Hash de la contraseña (Costo de encriptación estándar: 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return models.Usuario{}, err
	}

	nuevoUsuario := models.Usuario{
		Username: username,
		Password: string(hashedPassword), // Guardamos el hash, no el texto plano
		Rol:      rol,
	}

	err = s.repo.Crear(&nuevoUsuario)
	return nuevoUsuario, err
}

// Login compara el hash guardado contra la contraseña que ingresa el usuario
func (s *AuthService) Login(username, password string) (string, string, error) {
	usuario, err := s.repo.ObtenerPorUsername(username)
	if err != nil {
		return "", "", errors.New("credenciales incorrectas: usuario no existe")
	}

	// 🔐 Comparamos el hash de la Base de Datos con la contraseña ingresada
	err = bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("credenciales incorrectas: contraseña errónea")
	}

	// Generación del Token JWT
	claims := jwt.MapClaims{
		"username": usuario.Username,
		"rol":      usuario.Rol,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ClaveSecreta)
	if err != nil {
		return "", "", err
	}

	return tokenString, usuario.Rol, nil
}
