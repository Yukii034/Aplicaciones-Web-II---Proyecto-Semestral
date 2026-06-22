package modulo_pi // importar la carpeta en la que están

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service" //importar el service para los errores
	"proyecto-semestral/internal/storage"
)

type InventarioService struct {
	repo storage.InventarioRepository
}

func NewInventarioService(repo storage.InventarioRepository) *InventarioService {
	return &InventarioService{repo: repo}
}

func (s *InventarioService) ListarInventario() []models.Inventario {
	return s.repo.ListarInventario()
}

func (s *InventarioService) BuscarInventario(id int) (models.Inventario, error) {
	c, ok := s.repo.BuscarInventarioPorID(id)
	if !ok {
		return models.Inventario{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *InventarioService) CrearInventario(p models.Inventario) (models.Inventario, error) {
	if err := validarInventario(p); err != nil {
		return models.Inventario{}, err
	}

	return s.repo.CrearInventario(p), nil
}

func (s *InventarioService) ActualizarInventario(id int, p models.Inventario) (models.Inventario, error) {
	if err := validarInventario(p); err != nil {
		return models.Inventario{}, err
	}

	c, ok := s.repo.ActualizarInventario(id, p)
	if !ok {
		return models.Inventario{}, se.ErrNoEncontrado
	}

	return c, nil
}

func (s *InventarioService) BorrarInventario(id int) error {
	if !s.repo.BorrarInventario(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarInventario(p models.Inventario) error {
	if p.Nombre == "" {
		return se.ErrVacio
	}

	return nil
}
