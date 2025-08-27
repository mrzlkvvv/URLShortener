package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/mrzlkvvv/URLShortener/internal/database"
	"github.com/mrzlkvvv/URLShortener/internal/server/response"
)

func Redirect(urlGetter database.URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := zap.L()

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			render.JSON(w, r, response.Error("param 'alias' is empty"))
			l.Info("param 'alias' is empty")
			return
		}

		url, err := urlGetter.GetURL(r.Context(), alias)
		if err != nil {
			if errors.Is(err, database.ErrUrlIsNotExists) {
				render.JSON(w, r, response.Error("url not found"))
				l.Info("url not found", zap.String("alias", alias))
				return
			}

			render.JSON(w, r, response.Error("url getting is failed"))
			l.Info("URL getting is failed")
			return
		}

		l.Info("Success redirect", zap.String("alias", alias), zap.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)
	}
}
