package model

type Genre struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Slug string `gorm:"unique"`
}
