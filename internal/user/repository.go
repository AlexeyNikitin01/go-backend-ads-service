package user

import "context"
//go:generate mockery --output ../tests/mocks --name RepositoryUser
type RepositoryUser interface {
	UpdateUser(ctx context.Context, nickname string, email string, userID int64, activate bool) (*User, error)
	AddUser(ctx context.Context, user *User) (int64, error)
	CheckUser(ctx context.Context, user_id int64) (bool)
	GetUser(ctx context.Context, user_id int64) (*User, error)
	DeleteUser(ctx context.Context, user_id int64) (error)
}

//go:generate mockery --output ../tests/mocks --name RepositoryDbUser
type RepositoryDbUser interface {
	CreateUserDb(user UserDb) (int, error)
}
