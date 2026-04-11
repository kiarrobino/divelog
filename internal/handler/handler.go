package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kiarrobino/divelog/internal/exporter"
	"github.com/kiarrobino/divelog/internal/model"
	"github.com/kiarrobino/divelog/internal/service"
)

type DiveHandler struct {
	svc *service.DiveService
}

func NewDiveHandler(svc *service.DiveService) *DiveHandler {
	return &DiveHandler{svc: svc}
}

func (h *DiveHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input model.CreateDiveInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON"})
		return
	}

	newDive, err := h.svc.CreateDive(r.Context(), input)
	if err != nil {
		if errors.Is(err, model.ErrInvalidDepth) ||
			errors.Is(err, model.ErrInvalidDuration) ||
			errors.Is(err, model.ErrInvalidRating) ||
			errors.Is(err, model.ErrInvalidDate) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newDive)
}

func (h *DiveHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	dive, err := h.svc.GetDive(r.Context(), id)
	if err != nil {
		if errors.Is(err, model.ErrDiveNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dive)
}

func (h *DiveHandler) List(w http.ResponseWriter, r *http.Request) {
	limit := 20
	offset := 0
	var err error

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, err = strconv.Atoi(l)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		offset, err = strconv.Atoi(o)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	dives, err := h.svc.ListDives(r.Context(), limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dives)
}

func (h *DiveHandler) NDL(w http.ResponseWriter, r *http.Request) {
	var input model.NDLInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		if errors.Is(err, model.ErrInvalidDepth) ||
			errors.Is(err, model.ErrInvalidO2Percent) {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ndl, err := h.svc.CalculateNDL(input.Depth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ndl)
}

func (h *DiveHandler) Export(w http.ResponseWriter, r *http.Request) {
	dives, err := h.svc.ListDives(r.Context(), 10000, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", `attachment; filename="divelog.csv"`)

	err = exporter.WriteCSV(w, dives)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
