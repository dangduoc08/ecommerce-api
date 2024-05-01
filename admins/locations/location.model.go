package locations

type LocationModel struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name"`
	Slug       string         `json:"slug"`
	LocationID *uint          `json:"-" gorm:"nullable"`
	Location   *LocationModel `json:"location,omitempty" gorm:"references:ID;foreignKey:LocationID"`
}

type Query struct {
	LocationID *uint
}
