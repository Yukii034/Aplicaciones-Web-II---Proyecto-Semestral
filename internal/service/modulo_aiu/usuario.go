package modulo_aiu

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type UsuarioService struct {
	repo storage.UserRepository
}

func NewUsuarioService(repo storage.UserRepository) *UsuarioService {
	return &UsuarioService{repo: repo}
}

func (s *UsuarioService) ListarUsuarios() []models.Usuario {
	return s.repo.ListarUsuarios()
}

func (s *UsuarioService) BuscarUsuario(id int) (models.Usuario, error) {
	c, ok := s.repo.BuscarUsuarioPorID(id)
	if !ok {
		return models.Usuario{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *UsuarioService) CrearUsuario(u models.Usuario) (models.Usuario, error) {
	if err := validarUsuario(u); err != nil {
		return models.Usuario{}, err
	}
	creado := s.repo.CrearUsuario(u)
	if creado.ID == 0 {
		// El insert falló de verdad (ej. email duplicado -> viola unique),
		// aunque el repo no propague un error Go explícito.
		return models.Usuario{}, se.ErrRelacionInvalida
	}
	return creado, nil
}

func (s *UsuarioService) ActualizarUsuario(id int, u models.Usuario) (models.Usuario, error) {
	if err := validarUsuario(u); err != nil {
		return models.Usuario{}, err
	}
	c, ok := s.repo.ActualizarUsuario(id, u)
	if !ok {
		return models.Usuario{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *UsuarioService) BorrarUsuario(id int) error {
	if !s.repo.BorrarUsuario(id) {
		return se.ErrNoEncontrado
	}
	return nil
}

func validarUsuario(u models.Usuario) error {
	if u.Nombre == "" {
		return se.ErrVacio
	}
	if u.Email == "" {
		return se.ErrVacio
	}
	return nil
}
