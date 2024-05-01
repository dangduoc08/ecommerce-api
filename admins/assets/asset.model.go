package assets

import "time"

type AssetModel struct {
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	Extension string    `json:"extension"`
	IsDir     bool      `json:"is_dir"`
	UpdatedAt time.Time `json:"updated_at"`
}
