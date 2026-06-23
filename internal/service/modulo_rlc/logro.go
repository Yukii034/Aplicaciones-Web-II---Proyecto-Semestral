package modulo_rlc

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type LogroService struct {
	repo storage.LogroRepository
}

func NewLogroService(repo storage.LogroRepository) *LogroService {
	return &LogroService{repo: repo}
}

func (s *LogroService) ListarLogro() []models.Logro {
	return s.repo.ListarLogro()
}

func (s *LogroService) BuscarLogro(id int) (models.Logro, error) {
	c, ok := s.repo.BuscarLogroPorID(id)
	if !ok {
		return models.Logro{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *LogroService) CrearLogro(l models.Logro) (models.Logro, error) {
	if err := validarLogro(l); err != nil {
		return models.Logro{}, err
	}

	return s.repo.CrearLogro(l), nil
}

func (s *LogroService) ActualizarLogro(id int, l models.Logro) (models.Logro, error) {
	if err := validarLogro(l); err != nil {
		return models.Logro{}, err
	}

	c, ok := s.repo.ActualizarLogro(id, l)
	if !ok {
		return models.Logro{}, se.ErrNoEncontrado
	}

	return c, nil
}

func (s *LogroService) BorrarLogro(id int) error {
	if !s.repo.BorrarLogro(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarLogro(l models.Logro) error {
	if l.Nombre == "" {
		return se.ErrVacio
	}

	return nil
}
