package modulo_rlc

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type CalificacionService struct {
	repo    storage.CalificacionRepository
	usuario storage.UserRepository
	acuerdo storage.AcuerdoRepository
}

func NewCalificacionService(repo storage.CalificacionRepository, usuario storage.UserRepository, acuerdo storage.AcuerdoRepository) *CalificacionService {
	return &CalificacionService{repo: repo, usuario: usuario, acuerdo: acuerdo}
}

func (s *CalificacionService) ListarCalificacion() []models.Calificacion {
	return s.repo.ListarCalificacion()
}

func (s *CalificacionService) BuscarCalificacion(id int) (models.Calificacion, error) {
	c, ok := s.repo.BuscarCalificacionPorID(id)
	if !ok {
		return models.Calificacion{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *CalificacionService) CrearCalificacion(c models.Calificacion) (models.Calificacion, error) {
	if err := validarCalificacion(c); err != nil {
		return models.Calificacion{}, err
	}

	if _, ok := s.usuario.BuscarUsuarioPorID(c.UsuarioID); !ok {
		return models.Calificacion{}, se.ErrNoEncontrado
	}

	if _, ok := s.acuerdo.BuscarAcuerdoPorID(c.AcuerdoID); !ok {
		return models.Calificacion{}, se.ErrNoEncontrado
	}

	return s.repo.CrearCalificacion(c), nil
}

func (s *CalificacionService) ActualizarCalificacion(id int, c models.Calificacion) (models.Calificacion, error) {
	if err := validarCalificacion(c); err != nil {
		return models.Calificacion{}, err
	}

	c, ok := s.repo.ActualizarCalificacion(id, c)
	if !ok {
		return models.Calificacion{}, se.ErrNoEncontrado
	}

	return c, nil
}

func (s *CalificacionService) BorrarCalificacion(id int) error {
	if !s.repo.BorrarCalificacion(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarCalificacion(c models.Calificacion) error {
	if c.Comentarios == "" {
		return se.ErrVacio
	}
	if c.UsuarioID == 0 {
		return se.ErrVacio
	}
	return nil
}
