package service

import (
	"time"

	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/storage"

	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretJWT = []byte("cualquier_cosa_secreta")

const duracionToken = 24 * time.Hour

type Claims struct {
	Email     string `json:"email"`
	UsuarioID int    `json:"usuario_id"`
	jwt.RegisteredClaims
}

type AuthService struct {
	repo storage.UserRepository
}

func NewAuthService(repo storage.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Registrar(email, password string) (models.Usuario, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return models.Usuario{}, ErrCredencialesInvalidas
	}
	if _, existe := s.repo.BuscarUsuarioPorEmail(email); existe {
		return models.Usuario{}, ErrEmailEnUso
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.Usuario{}, err
	}
	return s.repo.CrearUsuario(models.Usuario{
		Email:        email,
		PasswordHash: string(hash),
	}), nil // agg luego errores
}

func (s *AuthService) Login(email, password string) (string, error) {
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	if email == "" || password == "" {
		return "", ErrCredencialesInvalidas
	}

	u, existe := s.repo.BuscarUsuarioPorEmail(email)
	if !existe {
		return "", ErrCredencialesInvalidas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrCredencialesInvalidas
	}
	return s.GenerarToken(u) //Linea reutilizable
}

func (s *AuthService) GenerarToken(u models.Usuario) (string, error) {
	claims := &Claims{
		UsuarioID: u.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duracionToken)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretJWT)
}

// Recoge el token, lo guarda en Claims y lo valida
func (s *AuthService) ValidarToken(token string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrCredencialesInvalidas
		}
		return secretJWT, nil
	})
	if err != nil || !parsedToken.Valid {
		return 0, ErrCredencialesInvalidas
	}
	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		return 0, ErrCredencialesInvalidas
	}
	return claims.UsuarioID, nil
}
