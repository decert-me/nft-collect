package global

import (
	"time"
)

type MODEL struct {
	ID        string    `gorm:"type:uuid;default:uuid_generate_v4();primarykey" json:"id"` // 主键
	CreatedAt time.Time `json:"-"`                                                         // 创建时间
	UpdatedAt time.Time `json:"-"`                                                         // 更新时间
}
