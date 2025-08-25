package storage

import "errors"

var (
	ErrUrlIsNotExists   = errors.New("URL is not exists")
	ErrUrlAlreadyExists = errors.New("URL already exists")
)

type Storage interface {
	URLSaver
	URLGetter
	Shutdown() error
}

type URLSaver interface {
	SaveURL(alias, url string) error
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}
