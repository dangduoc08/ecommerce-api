package models

type Location struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	LocationID *uint     `json:"-" gorm:"nullable"`
	Location   *Location `json:"location,omitempty" gorm:"references:ID;foreignKey:LocationID"`
}
