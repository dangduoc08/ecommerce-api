package user

type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"unique" json:"username"`
	Password  string `json:"-"`
	Salt      string `json:"-"`
	Email     string `gorm:"unique" json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
