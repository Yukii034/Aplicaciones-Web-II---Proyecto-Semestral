package modulo_aiu

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type AcuerdoItemService struct {
	repo storage.Acuerdo_ItemRepository
}

func NewAcuerdoItemService(repo storage.Acuerdo_ItemRepository) *AcuerdoItemService {
	return &AcuerdoItemService{repo: repo}
}

func (s *AcuerdoItemService) ListarAcuerdoItems() []models.AcuerdoItem {
	return s.repo.ListarAcuerdoItems()
}

func (s *AcuerdoItemService) BuscarAcuerdoItem(id int) (models.AcuerdoItem, error) {
	c, ok := s.repo.BuscarAcuerdoItemPorID(id)
	if !ok {
		return models.AcuerdoItem{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *AcuerdoItemService) CrearAcuerdoItem(a models.AcuerdoItem) (models.AcuerdoItem, error) {
	if err := validarAcuerdoItem(a); err != nil {
		return models.AcuerdoItem{}, err
	}
	return s.repo.CrearAcuerdoItem(a), nil
}

func (s *AcuerdoItemService) ActualizarAcuerdoItem(id int, a models.AcuerdoItem) (models.AcuerdoItem, error) {
	if err := validarAcuerdoItem(a); err != nil {
		return models.AcuerdoItem{}, err
	}
	c, ok := s.repo.ActualizarAcuerdoItem(id, a)
	if !ok {
		return models.AcuerdoItem{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *AcuerdoItemService) BorrarAcuerdoItem(id int) error {
	if !s.repo.BorrarAcuerdoItem(id) {
		return se.ErrNoEncontrado
	}
	return nil
}

func validarAcuerdoItem(a models.AcuerdoItem) error {
	if a.ID == 0 {
		return se.ErrVacio
	}
	return nil
}
