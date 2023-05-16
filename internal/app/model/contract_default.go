package model

import (
	"time"
)

type ContractDefault struct {
	ID         uint      `gorm:"primarykey"`
	CreatedAt  time.Time `json:"-"`                // 创建时间
	ContractID string    `gorm:"type:uuid;unique"` //
}
