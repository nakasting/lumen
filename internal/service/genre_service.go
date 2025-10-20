package service

import (
	"lumen/internal/dto"
	"lumen/internal/model"
	"lumen/internal/repository"
	"lumen/internal/utils"
)

type GenreService interface {
	FindAll() ([]dto.GenreRes, error)
	Create(genreReq *dto.GenreReq) (*dto.GenreRes, error)
	Exists(genreReq *dto.GenreReq, id uint) (bool, error)
}

type genreService struct {
	repo repository.GenreRepo
}

func NewGenreService(repo repository.GenreRepo) GenreService {
	return &genreService{repo: repo}
}

func (s *genreService) FindAll() ([]dto.GenreRes, error) {
	genres, err := s.repo.FindAll()

	if err != nil {
		return nil, err
	}

	var genreResponses []dto.GenreRes

	for _, g := range genres {
		genreResponses = append(genreResponses, dto.ToGenreRes(&g))
	}

	return genreResponses, nil
}

func (s *genreService) Create(genreReq *dto.GenreReq) (*dto.GenreRes, error) {
	slug := utils.GenerateSlug(genreReq.Name)

	newGenre := model.Genre{
		Name: genreReq.Name,
		Slug: slug,
	}

	if err := s.repo.Create(&newGenre); err != nil {
		return nil, err
	}

	genreRes := dto.ToGenreRes(&newGenre)

	return &genreRes, nil
}

func (s *genreService) Exists(genreReq *dto.GenreReq, id uint) (bool, error) {

	slug := utils.GenerateSlug(genreReq.Name)

	exists, err := s.repo.ExistsBySlug(slug, id)

	if err != nil {
		return false, err
	}

	return exists, nil
}
