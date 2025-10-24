package handler

import (
	"encoding/json"
	"fmt"
	"lumen/internal/dto"
	"lumen/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type GenreHandler struct {
	service   service.GenreService
	validator *validator.Validate
	logger    *zap.Logger
}

func NewGenreHandler(s service.GenreService, v *validator.Validate, l *zap.Logger) *GenreHandler {
	return &GenreHandler{
		service:   s,
		validator: v,
		logger:    l,
	}
}

// Routes
func (h *GenreHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.GetGenres)
	r.Post("/", h.CreateGenre)
	r.Get("/{id}", h.GetGenre)
	r.Put("/{id}", h.UpdateGenre)
}

/* ###################################################################################### */
func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	genreResponses, err := h.service.FindAll()

	if err != nil {
		h.logger.Error("Error getting genres", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(&genreResponses)
}

/* ###################################################################################### */
func (h *GenreHandler) GetGenre(w http.ResponseWriter, r *http.Request) {
	genreId := chi.URLParam(r, "id")
	enc := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")

	exists, err := h.service.ExistsByID(genreId)

	if err != nil {
		h.logger.Error("Genrename existing error", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
	}

	if !exists {
		h.logger.Warn("Genre not exists", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		errorResponse := dto.NewErrorResponse("Genre not found")
		enc.Encode(errorResponse)
		return
	}

	genreResponse, err := h.service.FindByID(genreId)

	if err != nil {
		h.logger.Error("Cannot get genre", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(genreResponse)
}

/* ###################################################################################### */
func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	var genreRequest dto.GenreReq
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)

	if err := dec.Decode(&genreRequest); err != nil {
		h.logger.Error("Error encoding json", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	if err := h.validator.Struct(genreRequest); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		var fieldErrors []dto.FieldError
		validationErrors := err.(validator.ValidationErrors)

		for _, validationError := range validationErrors {
			fieldError := dto.NewFieldError(validationError.Field(), validationError.Error())
			fieldErrors = append(fieldErrors, *fieldError)
		}
		errorRes := dto.NewValidationErrors(fieldErrors)
		enc.Encode(&errorRes)
		return
	}

	genreExists, err := h.service.Exists(&genreRequest, "")

	if err != nil {
		h.logger.Error("Error checking on database", zap.Error(err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	if genreExists {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		var fieldErrors []dto.FieldError
		fieldError := dto.NewFieldError("Name", fmt.Sprintf("Genre with name '%s' is already exists", genreRequest.Name))
		fieldErrors = append(fieldErrors, *fieldError)
		errorRes := dto.NewValidationErrors(fieldErrors)
		enc.Encode(&errorRes)
		return
	}

	genreResponse, err := h.service.Create(&genreRequest)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(&errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	enc.Encode(genreResponse)
}

func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	genreId := chi.URLParam(r, "id")
	dec := json.NewDecoder(r.Body)
	enc := json.NewEncoder(w)

	exists, err := h.service.ExistsByID(genreId)

	if err != nil {
		h.logger.Error("Error exists by ID", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	if !exists {
		h.logger.Warn("Genre not exists", zap.Error(err))
		w.WriteHeader(http.StatusNotFound)
		errorResponse := dto.NewErrorResponse("Genre not found")
		enc.Encode(errorResponse)
		return
	}

	var genreRequest dto.GenreReq

	if err := dec.Decode(&genreRequest); err != nil {
		h.logger.Error("Cannot decode request", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	if err := h.validator.Struct(genreRequest); err != nil {
		h.logger.Warn("Error validate genre", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)

		var fieldErrors []dto.FieldError
		validationErrors := err.(validator.ValidationErrors)

		for _, validationError := range validationErrors {
			fieldError := dto.NewFieldError(validationError.Field(), validationError.Error())
			fieldErrors = append(fieldErrors, *fieldError)
		}
		errorRes := dto.NewValidationErrors(fieldErrors)
		enc.Encode(&errorRes)
		return
	}

	genreResponse, err := h.service.Update(genreId, &genreRequest)

	if err != nil {
		h.logger.Error("Cannot update genre", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		errorResponse := dto.NewErrorResponse("Error internal server")
		enc.Encode(errorResponse)
		return
	}

	w.WriteHeader(http.StatusOK)
	enc.Encode(genreResponse)
}
