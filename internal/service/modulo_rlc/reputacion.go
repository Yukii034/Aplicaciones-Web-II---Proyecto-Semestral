package modulo_rlc

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type ReputacionService struct {
	repo    storage.ReputacionRepository
	usuario storage.UserRepository
}

func NewReputacionService(repo storage.ReputacionRepository) *ReputacionService {
	return &ReputacionService{repo: repo}
}

func (s *ReputacionService) ListarReputacion() []models.Reputacion {
	return s.repo.ListarReputacion()
}

func (s *ReputacionService) BuscarReputacion(id int) (models.Reputacion, error) {
	c, ok := s.repo.BuscarReputacionPorID(id)
	if !ok {
		return models.Reputacion{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *ReputacionService) CrearReputacion(r models.Reputacion) (models.Reputacion, error) {
	if err := validarReputacion(r); err != nil {
		return models.Reputacion{}, err
	}

	if _, ok := s.usuario.BuscarUsuarioPorID(r.UsuarioID); !ok {
		return models.Reputacion{}, se.ErrNoEncontrado
	}

	return s.repo.CrearReputacion(r), nil
}

func (s *ReputacionService) ActualizarReputacion(id int, r models.Reputacion) (models.Reputacion, error) {
	if err := validarReputacion(r); err != nil {
		return models.Reputacion{}, err
	}

	c, ok := s.repo.ActualizarReputacion(id, r)
	if !ok {
		return models.Reputacion{}, se.ErrNoEncontrado
	}

	return c, nil
}

func (s *ReputacionService) BorrarReputacion(id int) error {
	if !s.repo.BorrarReputacion(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarReputacion(p models.Reputacion) error {
	if p.PuntosTotales == 0 {
		return se.ErrVacio
	}
	if p.UsuarioID == 0 {
		return se.ErrVacio
	}
	return nil
}
