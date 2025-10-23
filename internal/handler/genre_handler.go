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
	r.Post("/", h.CreateGenre)
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

	genreExists, err := h.service.Exists(&genreRequest, 0)

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
