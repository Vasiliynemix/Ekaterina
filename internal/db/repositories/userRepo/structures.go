package userRepo

type UserAddParams struct {
	TelegramID int64
	UserName   string
	CreatedAt  int64
	IsAdmin    bool
}

type UserShow struct {
	TelegramID int64
	UserName   string
	IsAdmin    bool
	IsModer    bool
}
