package service

import "errors"

var (
	ErrCredencialesInvalidas = errors.New("email o contraseña incorrectos")
	ErrEmailEnUso            = errors.New("email ya en uso")
	ErrNoEncontrado          = errors.New("No encontrado")
	ErrVacio                 = errors.New("Campo vacío")
	ErrRelacionInvalida      = errors.New("Uno o más IDs relacionados no existen (violación de llave foránea)")
)
