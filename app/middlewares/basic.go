package middlewares

import (
	"net/http"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-skeleton-chi/app/configs"
)

func SetAuthBasic() func(http.Handler) http.Handler {
	cfg := configs.GetInstance()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Basic Authentication credentials
			user, password, hasAuth := r.BasicAuth()
			isValid := hasAuth && user == cfg.App.BasicUsername && password == cfg.App.BasicPassword
			if !isValid {
				err := response_mapper.NewError(response_mapper.ErrUnauthorized, response_mapper.NewResponseMultiLang(
					response_mapper.MultiLanguages{
						ID: "Token tidak valid",
						EN: "Invalid token",
					},
				))

				response_mapper.RenderJSON(w, http.StatusUnauthorized, err)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
