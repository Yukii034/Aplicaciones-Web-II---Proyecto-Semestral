package modulo_pi

import (
	"proyecto-semestral/internal/models"
	se "proyecto-semestral/internal/service"
	"proyecto-semestral/internal/storage"
)

type PublicacionService struct {
	repo       storage.PublicacionRepository
	inventario storage.InventarioRepository // para las relaciones
	usuario    storage.UserRepository
}

func NewPublicacionService(repo storage.PublicacionRepository, inv storage.InventarioRepository, usr storage.UserRepository) *PublicacionService {
	return &PublicacionService{repo: repo, inventario: inv, usuario: usr}
}

func (s *PublicacionService) ListarPublicacion() []models.Publicacion {
	return s.repo.ListarPublicacion()
}

func (s *PublicacionService) BuscarPublicacion(id int) (models.Publicacion, error) {
	c, ok := s.repo.BuscarPublicacionPorID(id)
	if !ok {
		return models.Publicacion{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *PublicacionService) CrearPublicacion(p models.Publicacion) (models.Publicacion, error) {
	if err := validarPublicacion(p); err != nil {
		return models.Publicacion{}, err
	}

	// verifica que el inventario referenciado exista
	if _, ok := s.inventario.BuscarInventarioPorID(p.InventarioID); !ok {
		return models.Publicacion{}, se.ErrNoEncontrado
	}

	// verifica que el usuario referenciado exista
	if _, ok := s.usuario.BuscarUsuarioPorID(p.UsuarioID); !ok {
		return models.Publicacion{}, se.ErrNoEncontrado
	}

	return s.repo.CrearPublicacion(p), nil
}

func (s *PublicacionService) ActualizarPublicacion(id int, p models.Publicacion) (models.Publicacion, error) {
	if err := validarPublicacion(p); err != nil {
		return models.Publicacion{}, err
	}

	// verifica que el inventario referenciado exista
	if _, ok := s.inventario.BuscarInventarioPorID(p.InventarioID); !ok {
		return models.Publicacion{}, se.ErrNoEncontrado
	}

	c, ok := s.repo.ActualizarPublicacion(id, p)
	if !ok {
		return models.Publicacion{}, se.ErrNoEncontrado
	}
	return c, nil
}

func (s *PublicacionService) BorrarPublicacion(id int) error {
	if !s.repo.BorrarPublicacion(id) {
		return se.ErrNoEncontrado
	}

	return nil
}

func validarPublicacion(p models.Publicacion) error {
	if p.Titulo == "" {
		return se.ErrVacio
	}

	return nil
}
