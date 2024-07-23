package router

import (
	"net/http"
	"time"

	response_mapper "github.com/adamnasrudin03/go-helpers/response-mapper/v1"
	"github.com/adamnasrudin03/go-skeleton-chi/app/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

type routes struct {
	HttpServer *chi.Mux
}

func NewRoutes(h controller.Controllers) routes {
	var err error
	r := routes{
		HttpServer: chi.NewRouter(),
	}

	r.HttpServer.Use(middleware.Logger)
	r.HttpServer.Use(middleware.Recoverer)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.HttpServer.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.HttpServer.Use(middleware.Timeout(60 * time.Second))
	r.HttpServer.Use(render.SetContentType(render.ContentTypeJSON))

	r.HttpServer.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response_mapper.RenderJSON(w, http.StatusOK, response_mapper.NewResponseMultiLang(response_mapper.MultiLanguages{
			ID: "selamat datang di server ini",
			EN: "Welcome this server",
		}))
	})

	r.HttpServer.NotFound(func(w http.ResponseWriter, r *http.Request) {
		err = response_mapper.ErrRouteNotFound()
		response_mapper.RenderJSON(w, http.StatusNotFound, err)
	})
	return r
}

func (r routes) Run(addr string) error {
	server := &http.Server{Addr: addr, Handler: r.HttpServer}
	return server.ListenAndServe()
}
