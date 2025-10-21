package handler

import (
	"encoding/json"
	"fmt"
	"lumen/internal/dto"
	"lumen/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type Genrehandler struct {
	service   service.GenreService
	validator *validator.Validate
	logger    *zap.Logger
}

func NewGenreHandler(s service.GenreService, v *validator.Validate, l *zap.Logger) *Genrehandler {
	return &Genrehandler{
		service:   s,
		validator: v,
		logger:    l,
	}
}

// Routes
func (h *Genrehandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetGenres)
	r.Post("/", h.CreateGenre)
	r.Put("/{id}", h.UpdateGenre)
}

/* ###################################################################################### */
func (h *Genrehandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	genres, err := h.service.FindAll()

	if err != nil {
		h.logger.Error("Cannot get genres on database", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "body",
					Error: "Error internal server",
				},
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genres)

}

/* ###################################################################################### */
func (h *Genrehandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var genreReq dto.GenreReq

	if err := json.NewDecoder(r.Body).Decode(&genreReq); err != nil {
		h.logger.Warn("Error decode request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorResponse{Message: "Error internal server"})
		return
	}

	if err := h.validator.Struct(genreReq); err != nil {
		h.logger.Warn("Error validate model", zap.Error(err))
		errs := err.(validator.ValidationErrors)
		var errFields []dto.ErrorField
		for _, e := range errs {
			errFields = append(errFields, dto.ErrorField{
				Field: strings.ToLower(e.Field()),
				Error: e.Error(),
			})
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: errFields,
		})
		return
	}

	exists, err := h.service.Exists(&genreReq, 0)

	if err != nil {
		h.logger.Error("Cannot validate genre exists", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "body",
					Error: "Error internal server",
				},
			},
		})
	}

	if exists {
		h.logger.Warn("Genre already exists", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "name",
					Error: fmt.Sprintf("Genre with name %s is already exists", genreReq.Name),
				},
			},
		})
		return
	}

	res, err := h.service.Create(&genreReq)

	if err != nil {
		h.logger.Error("Cannot save model", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "body",
					Error: "Error internal server",
				},
			},
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

/* ###################################################################################### */
func (h *Genrehandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	var genreReq dto.GenreReq

	if err := json.NewDecoder(r.Body).Decode(&genreReq); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "body",
					Error: "Cannot decode request",
				},
			},
		})
		return
	}

	genreId, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("Invalid param", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "url",
					Error: "Invalid param",
				},
			},
		})
		return
	}

	genreRes, err := h.service.Update(uint(genreId), &genreReq)

	if err != nil {
		h.logger.Warn("Cannot update model", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.ErrorRes{
			Errors: []dto.ErrorField{
				dto.ErrorField{
					Field: "body",
					Error: "Error internal server",
				},
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(genreRes)
}
