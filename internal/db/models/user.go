package models

type User struct {
	Base
	UserName   string `gorm:"unique"`
	TelegramID int64  `gorm:"unique"`
	IsAdmin    bool   `gorm:"default:false"`
	IsModer    bool   `gorm:"default:false"`

	TypeSchedule string

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}
