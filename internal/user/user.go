package user

type User struct {
	UserID  int64
	NickName string
	Email   string
	Activate bool
}

type UserDb struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
