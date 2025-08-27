package database

import (
	"context"
	"errors"
)

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
	SaveURL(ctx context.Context, alias, url string) error
}

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}
