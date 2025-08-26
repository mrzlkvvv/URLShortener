package database

import "errors"

var (
	ErrUrlIsNotExists   = errors.New("URL is not exists")
	ErrUrlAlreadyExists = errors.New("URL already exists")
)

type Database interface {
	URLSaver
	URLGetter
	Shutdown()
}

type URLSaver interface {
	SaveURL(alias, url string) error
}

type URLGetter interface {
	GetURL(alias string) (string, error)
}
