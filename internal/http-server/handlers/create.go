package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/http-server/response"
	"github.com/mrzlkvvv/URLShortener/internal/logger"
	"github.com/mrzlkvvv/URLShortener/internal/random"
	"github.com/mrzlkvvv/URLShortener/internal/storage"
)

const ALIAS_LENGTH = 6

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias"`
}

type URLSaver interface {
	SaveURL(alias, url string) error
}

func Create(urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logger.New()

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error("failed to decode request"))
			l.Error("Request decoding is failed", zap.Error(err))
			return
		}

		err = validator.New().Struct(req)
		if err != nil {
			validationErrs := err.(validator.ValidationErrors)
			render.JSON(w, r, response.ValidationError(validationErrs))
			l.Error("Invalid request", zap.Error(err))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(ALIAS_LENGTH)
		}

		err = urlSaver.SaveURL(alias, req.URL)
		if err != nil {

			if errors.Is(err, storage.ErrUrlAlreadyExists) {
				render.JSON(w, r, response.Error("url already exists"))
				l.Error("URL saving is failed", zap.Error(err))
				return
			}

			render.JSON(w, r, response.Error("url saving is failed"))
			l.Error("URL saving is failed", zap.Error(err))
			return
		}

		render.JSON(w, r, Response{Response: response.OK(), Alias: alias})
		l.Info("URL was saved", zap.String("alias", alias), zap.String("url", req.URL))
	}
}
