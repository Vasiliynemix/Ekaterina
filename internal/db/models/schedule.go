package models

type Schedule struct {
	Base

	TelegramID int64

	WeekEven WeekEven `gorm:"foreignKey:ScheduleID"`
	WeekOdd  WeekOdd  `gorm:"foreignKey:ScheduleID"`

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}

type WeekEven struct {
	Base

	ScheduleID uint

	Days []Day `gorm:"foreignKey:WeekEvenID"`

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}

type WeekOdd struct {
	Base

	ScheduleID uint

	Days []Day `gorm:"foreignKey:WeekOddID"`

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}

type Day struct {
	Base

	DayName string

	WeekEvenID uint
	WeekOddID  uint

	CreatedAt int64 `gorm:"not null"`
	UpdatedAt int64
}
