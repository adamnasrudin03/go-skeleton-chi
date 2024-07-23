package router

import (
	"github.com/adamnasrudin03/go-skeleton-chi/app/controller"
	"github.com/adamnasrudin03/go-skeleton-chi/app/middlewares"
	"github.com/go-chi/chi/v5"
)

func (r routes) teamMember(handler controller.TeamMemberController) chi.Router {

	r.router.Route("/v1/team-members", func(r chi.Router) {
		r.With(middlewares.SetAuthBasic()).Route("/", func(r chi.Router) {
			r.Post("/", handler.Create)
			r.Delete("/{id}", handler.Delete)
			r.Put("/{id}", handler.Update)
		})

		r.Get("/", handler.GetList)
		r.Get("/{id}", handler.GetDetail)
	})

	return r.router
}
