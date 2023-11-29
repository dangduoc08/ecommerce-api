package user

import "time"

type UserStatus string

const (
	ACTIVE    = "active"
	INACTIVE  = "inactive"
	SUSPENDED = "suspended"
)

type User struct {
	ID        uint       `gorm:"primaryKey;->" json:"id"`
	Username  string     `gorm:"unique;not null" json:"username"`
	Hash      string     `gorm:"not null" json:"-"`
	Email     string     `gorm:"unique;not null" json:"email"`
	FirstName string     `gorm:"not null" json:"first_name"`
	LastName  string     `gorm:"not null" json:"last_name"`
	Status    UserStatus `gorm:"not null;type:user_status" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime:true" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime:true" json:"updated_at"`
}
