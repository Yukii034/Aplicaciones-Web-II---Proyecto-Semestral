package modulo_aiu

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type AcuerdoService struct {
	repo storage.AcuerdoRepository
}

func NewAcuerdoService(repo storage.AcuerdoRepository) *AcuerdoService {
	return &AcuerdoService{repo: repo}
}

func (s *AcuerdoService) ListarAcuerdos() []models.Acuerdo {
	return s.repo.ListarAcuerdos()
}

func (s *AcuerdoService) BuscarAcuerdo(id int) (models.Acuerdo, error) {
	c, ok := s.repo.BuscarAcuerdoPorID(id)
	if !ok {
		return models.Acuerdo{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *AcuerdoService) CrearAcuerdo(a models.Acuerdo) (models.Acuerdo, error) {
	if err := validarAcuerdo(a); err != nil {
		return models.Acuerdo{}, err
	}
	return s.repo.CrearAcuerdo(a), nil
}

func (s *AcuerdoService) ActualizarAcuerdo(id int, a models.Acuerdo) (models.Acuerdo, error) {
	if err := validarAcuerdo(a); err != nil {
		return models.Acuerdo{}, err
	}
	c, ok := s.repo.ActualizarAcuerdo(id, a)
	if !ok {
		return models.Acuerdo{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *AcuerdoService) BorrarAcuerdo(id int) error {
	if !s.repo.BorrarAcuerdo(id) {
		return se.ErrNoEncontrado
	}
	return nil
}

func validarAcuerdo(a models.Acuerdo) error {
	// Valida campos reales de negocio, no el ID autoincremental
	if a.PublicacionID == 0 || a.Tipo == "" {
		return se.ErrVacio
	}
	return nil
}
