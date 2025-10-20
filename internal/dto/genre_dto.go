package dto

import "lumen/internal/model"

type GenreReq struct {
	Name string `json:"name" validate:"required,min=3,max=20,alphaunicode"`
}

type GenreRes struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func ToGenreRes(g *model.Genre) GenreRes {
	return GenreRes{
		ID:   g.ID,
		Name: g.Name,
		Slug: g.Slug,
	}
}
