package storage

import "errors"

var (
	ErrUrlIsNotExists   = errors.New("URL is not exists")
	ErrUrlAlreadyExists = errors.New("URL already exists")
)
