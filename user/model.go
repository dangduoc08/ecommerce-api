package user

type UserStatus string

const (
	ACTIVE    = "ACTIVE"
	INACTIVE  = "INACTIVE"
	SUSPENDED = "SUSPENDED"
)

type User struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Username  string     `gorm:"unique;not null" json:"username"`
	Hash      string     `gorm:"not null" json:"-"`
	Email     string     `gorm:"unique;not null" json:"email"`
	FirstName string     `gorm:"not null" json:"first_name"`
	LastName  string     `gorm:"not null" json:"last_name"`
	Status    UserStatus `gorm:"not null;type:user_status" json:"status"`
}
