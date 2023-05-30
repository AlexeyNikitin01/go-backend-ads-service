package userrepo

import (
	"context"
	"fmt"
	"sync"
	
	"ads/internal/user"
)

type keyUserId int64
type userStructType = *user.User

type UserRepositoryMap struct {
	sync.Mutex
	countUserID int64
	mapUser map[keyUserId]userStructType
}

func (ur *UserRepositoryMap) AddUser(ctx context.Context, user *user.User) (int64, error) {
	ur.countUserID += 1
	user.UserID = ur.countUserID
	ur.mapUser[keyUserId(ur.countUserID)] = user
	return ur.countUserID, nil
}

func (ur *UserRepositoryMap) UpdateUser(ctx context.Context, nickname string, email string, userId int64, activate bool) (*user.User, error) {
	user, ok := ur.mapUser[keyUserId(userId)]
	if !ok {
		return nil, fmt.Errorf("not user in map")
	}
	user.NickName = nickname
	user.Email = email
	user.Activate = activate

	return user, nil
}

func (ur *UserRepositoryMap) CheckUser(ctx context.Context, user_id int64) bool {	
	for _, user := range ur.mapUser {
		if user.UserID == int64(user_id) {
			return true
		}
	}
	return false
}

func (ur *UserRepositoryMap) GetUser(ctx context.Context, user_id int64) (*user.User, error) {
	user, ok := ur.mapUser[keyUserId(user_id)]

	if !ok {
		return nil, fmt.Errorf("not found user in db")
	}
	
	return user, nil
}

func (ur *UserRepositoryMap) DeleteUser(ctx context.Context, user_id int64) (error) {
	_, ok := ur.mapUser[keyUserId(user_id)]

	if !ok {
		return fmt.Errorf("not found in db")
	}
	
	delete(ur.mapUser, keyUserId(user_id))
	
	_, ok = ur.mapUser[keyUserId(user_id)]
	if ok {
		return fmt.Errorf("didn`t deleted user")
	}
	
	return nil
}

func New() user.RepositoryUser {
	return &UserRepositoryMap{
		countUserID: -1,
		mapUser: make(map[keyUserId]userStructType),
	}
}
