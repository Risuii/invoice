package v1

import (
	"net/http"

	"github.com/Risuii/invoice/src/v1/handler"
	"github.com/go-chi/chi/v5"
)

func Router(r *chi.Mux, deps *Dependency) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// product

	r.Route("/invoice/v1", func(v1 chi.Router) {
		v1.Post("/", handler.CreateInvoiceHandler(deps.Services.Invoicesvc))
		v1.Patch("/{id}", handler.UpdateInvoiceHandler(deps.Services.Invoicesvc))
		v1.Get("/", handler.GetListInvoicesHandler(deps.Services.Invoicesvc))
		v1.Get("/{id}", handler.GetDetailInvoicesHandler(deps.Services.Invoicesvc))
	})
}
