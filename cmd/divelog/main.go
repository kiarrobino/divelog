package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kiarrobino/divelog/internal/handler"
	"github.com/kiarrobino/divelog/internal/repository"
	"github.com/kiarrobino/divelog/internal/service"
)

func main() {
	repo, err := repository.NewSQLiteRepository("divelog.db")
	if err != nil {
		log.Fatal(err)
	}

	svc := service.NewDiveService(repo)
	h := handler.NewDiveHandler(svc)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/dives", h.Create)
	r.Get("/dives/{id}", h.GetByID)
	r.Get("/dives", h.List)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
