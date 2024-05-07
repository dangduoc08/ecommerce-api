package mailConfigurations

import "time"

type MailConfigurationModel struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Host      string    `json:"host"`
	Port      int       `json:"post"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	StoreID   uint      `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type Query struct {
	ID      uint
	StoreID uint
	Sort    string
	Order   string
	Limit   int
	Offset  int
}

type Update struct {
	ID       uint
	Host     string
	Port     int
	Username string
	Password string
}
