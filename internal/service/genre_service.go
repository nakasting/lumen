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
	FindByID(id string) (*dto.GenreRes, error)
	Update(id string, genreReq *dto.GenreReq) (*dto.GenreRes, error)
	Exists(genreReq *dto.GenreReq, id string) (bool, error)
	ExistsByID(id string) (bool, error)
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

func (s *genreService) FindByID(id string) (*dto.GenreRes, error) {
	genre, err := s.repo.FindByID(id)

	if err != nil {
		return nil, err
	}

	genreResponse := dto.ToGenreRes(genre)
	return &genreResponse, nil
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

func (s *genreService) Update(id string, genreReq *dto.GenreReq) (*dto.GenreRes, error) {
	genre, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	slug := utils.GenerateSlug(genreReq.Name)
	genre.Name = genreReq.Name
	genre.Slug = slug

	if err := s.repo.Update(genre); err != nil {
		return nil, err
	}

	genreResponse := dto.ToGenreRes(genre)
	return &genreResponse, nil
}

func (s *genreService) Exists(genreReq *dto.GenreReq, id string) (bool, error) {

	slug := utils.GenerateSlug(genreReq.Name)

	exists, err := s.repo.ExistsBySlug(slug, id)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *genreService) ExistsByID(id string) (bool, error) {
	exists, err := s.repo.ExistsByID(id)

	if err != nil {
		return false, err
	}

	return exists, nil
}
