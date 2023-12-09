package locations

type Location struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name"`
	Slug       string    `json:"slug"`
	LocationID *uint     `json:"location_id" gorm:"nullable"`
	Location   *Location `json:"location,omitempty" gorm:"references:ID;foreignKey:LocationID"`
}
