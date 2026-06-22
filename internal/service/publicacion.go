package service

import (
	"proyecto-semestral/internal/models"
	"proyecto-semestral/internal/storage"
)

type PublicacionService struct {
	repo storage.PublicacionRepository
}

func NewCategoriaService(repo storage.PublicacionRepository) *PublicacionService {
	return &PublicacionService{repo: repo}
}

func (s *PublicacionService) Listar() []models.Publicacion {
	return s.repo.ListarPublicacion()
}

func (s *PublicacionService) Obtener(id int) (models.Publicacion, error) {
	c, ok := s.repo.BuscarPublicacionPorID(id)
	if !ok {
		return models.Publicacion{}, ErrNoEncontrado
	}
	return c, nil
}

func (s *PublicacionService) Crear(p models.Publicacion) (models.Publicacion, error) {
	if err := validarPublicacion(p); err != nil {
		return models.Publicacion{}, err
	}

	return s.repo.CrearPublicacion(p), nil
}

func (s *PublicacionService) Actualizar(id int, p models.Publicacion) (models.Publicacion, error) {
	if err := validarPublicacion(p); err != nil {
		return models.Publicacion{}, err
	}

	c, ok := s.repo.ActualizarPublicacion(id, p)
	if !ok {
		return models.Publicacion{}, ErrNoEncontrado
	}

	return c, nil
}

func (s *PublicacionService) Borrar(id int) error {
	if !s.repo.BorrarPublicacion(id) {
		return ErrNoEncontrado
	}

	return nil
}

func validarPublicacion(p models.Publicacion) error {
	if p.Titulo == "" {
		return ErrVacio
	}

	return nil
}
