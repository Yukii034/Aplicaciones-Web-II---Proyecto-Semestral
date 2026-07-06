package modulo_rlc

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type Logro_UsuarioService struct {
	repo  storage.Logro_UsuarioRepository
	logro storage.LogroRepository
}

func NewLogro_UsuarioService(repo storage.Logro_UsuarioRepository) *Logro_UsuarioService {
	return &Logro_UsuarioService{repo: repo}
}

func (s *Logro_UsuarioService) ListarLogro_Usuario() []models.Logro_Usuario {
	return s.repo.ListarLogro_Usuario()
}

func (s *Logro_UsuarioService) BuscarLogro_Usuario(id int) (models.Logro_Usuario, error) {
	c, ok := s.repo.BuscarLogro_UsuarioPorID(id)
	if !ok {
		return models.Logro_Usuario{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *Logro_UsuarioService) CrearLogro_Usuario(lu models.Logro_Usuario) (models.Logro_Usuario, error) {
	if err := validarLogro_Usuario(lu); err != nil {
		return models.Logro_Usuario{}, err
	}

	return s.repo.CrearLogro_Usuario(lu), nil
}

func (s *Logro_UsuarioService) ActualizarLogro_Usuario(id int, lu models.Logro_Usuario) (models.Logro_Usuario, error) {
	if err := validarLogro_Usuario(lu); err != nil {
		return models.Logro_Usuario{}, err
	}

	c, ok := s.repo.ActualizarLogro_Usuario(id, lu)
	if !ok {
		return models.Logro_Usuario{}, se.ErrNoEncontrado
	}

	return c, nil
}

func (s *Logro_UsuarioService) BorrarLogro_Usuario(id int) error {
	if !s.repo.BorrarLogro_Usuario(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarLogro_Usuario(lu models.Logro_Usuario) error {
	if lu.FechaDesbl.IsZero() {
		return se.ErrVacio
	}

	return nil
}
