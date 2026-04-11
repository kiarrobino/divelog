package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kiarrobino/divelog/internal/config"
	"github.com/kiarrobino/divelog/internal/handler"
	"github.com/kiarrobino/divelog/internal/repository"
	"github.com/kiarrobino/divelog/internal/service"
)

func main() {
	conf := config.Load()

	repo, err := repository.NewSQLiteRepository(conf.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewDiveService(repo)
	h := handler.NewDiveHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(handler.MetricsMiddleware)

	r.Route("/api", func(r chi.Router) {
		r.Post("/dives", h.Create)
		r.Post("/ndl", h.NDL)
		r.Get("/health", h.Health)
		r.Get("/dives/{id}", h.GetByID)
		r.Get("/dives", h.List)
		r.Get("/export/csv", h.Export)
		r.Handle("/metrics", promhttp.Handler())
	})

	r.Handle("/*", http.FileServer(http.Dir("web/static")))

	log.Printf("listening on %s", conf.Addr)
	log.Fatal(http.ListenAndServe(conf.Addr, r))
}
