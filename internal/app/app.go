package app

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"ads/internal/ads"
	"ads/internal/user"

	"github.com/AlexeyNikitin01/validate"
)

var ErrBadRequest = fmt.Errorf("bad request")
var ErrForbidden = fmt.Errorf("forbidden")
var ErrNotFound = fmt.Errorf("not found user in db")

type validateStruct struct {
	Text string `json:"text"`
	Title string `json:"title"`
}
//go:generate mockery --output ../tests/mocks --name App
type App interface {
	AdApp
	UserApp
	UserDbApp
}

type appStruct struct {
	adApp
	userApp
	authApp
}

type AdApp interface {
	CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, adID int64, published bool, authorID int64) (*ads.Ad, error)
	UpdateAd(ctx context.Context, authorID int64, title string, text string, adID int64) (*ads.Ad, error)
	GetAd(ctx context.Context, adID int64) (*ads.Ad, error)
	ListAds(ctx context.Context) ([]*ads.Ad, error)
	SearchAdByName(ctx context.Context, title string) ([]*ads.Ad, error)
	ListAdsAuthor(ctx context.Context, author int64) ([]*ads.Ad, error)
	ListAdsDate(ctx context.Context, day int64) ([]*ads.Ad, error)
	DeleteAd(ctx context.Context, authorID int64, adID int64) (*ads.Ad, error)
}

type adApp struct {
	repository ads.RepositryAd
}

func (a *adApp) CreateAd(ctx context.Context, title string, text string, authorID int64) (*ads.Ad, error) {
	if err := validate.Validate(validateStruct{Text: text, Title: title}); err != nil {
		return nil, ErrBadRequest
	}
	
	ad := ads.Ad{Title: title, Text: text, AuthorID: authorID, Published: false, CreateDate: time.Now().UTC()}
	id, err := a.repository.Add(ctx, &ad)

	if err != nil {
		return nil, err
	}

	ad.ID = id

	return &ad, nil
}

func (a *adApp) ChangeAdStatus(ctx context.Context, adID int64, published bool, authorID int64) (*ads.Ad, error) {
	ad, err := a.repository.GetAd(ctx, adID)

	if err != nil {
		return nil, err
	} else if ad.AuthorID != authorID || ad.ID != adID {
		return nil, ErrForbidden
	}
	
	ad, err = a.repository.ChangeStatus(ctx, adID, published, authorID)
	
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a *adApp) UpdateAd(ctx context.Context, authorID int64, title string, text string, adID int64) (*ads.Ad, error) {
	if err := validate.Validate(validateStruct{Text: text, Title: title}); err != nil {
		return nil, ErrBadRequest
	}
	
	ad, err := a.repository.GetAd(ctx, adID)

	if err != nil {
		return nil, err
	} else if ad.AuthorID != authorID || ad.ID != adID {
		return nil, ErrForbidden
	}

	ad, err = a.repository.Update(ctx, authorID, title, text, adID)
	
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a *adApp) GetAd(ctx context.Context, adID int64) (*ads.Ad, error) {
	ad, err := a.repository.GetAd(ctx, adID)
	if err != nil {
		return nil, ErrBadRequest
	}
	return ad, nil
}

func (a *adApp) ListAds(ctx context.Context) ([]*ads.Ad, error) {
	ads, err := a.repository.ListAds(ctx)
	if err != nil {
		return nil, ErrBadRequest
	}
	return ads, nil
}

func (a *adApp) SearchAdByName(ctx context.Context, title string) ([]*ads.Ad, error) {
	ads, err := a.repository.Search(ctx, title)
	if err != nil {
		return nil, err
	}
	return ads, nil
}

func (a *adApp) ListAdsAuthor(ctx context.Context, author int64) ([]*ads.Ad, error) {
	ads, err := a.repository.ListAdsAuthor(ctx, author)
	if err != nil {
		return nil, ErrBadRequest
	}
	return ads, nil
}

func (a *adApp) ListAdsDate(ctx context.Context, day int64) ([]*ads.Ad, error) {
	ads, err := a.repository.ListAdsDate(ctx, day)
	if err != nil {
		return nil, ErrBadRequest
	}
	return ads, nil
}

func (a *adApp) DeleteAd(ctx context.Context, authorID int64, adID int64) (*ads.Ad, error) {
	ad, err := a.repository.DeleteAd(ctx, authorID, adID)
	if err != nil {
		return nil, err
	}
	return ad, nil
}
 
type UserApp interface {
	CreateUser(ctx context.Context, nickname string, email string) (*user.User, error)
	UpdateUser(ctx context.Context, nickname string, email string, userID int64, activate bool) (*user.User, error)
	CheckUser(ctx context.Context, userID int64) (error)
	GetUser(ctx context.Context, userID int64) (*user.User, error)
	DeleteUser(ctx context.Context, userID int64) (error)
}

type userApp struct {
	repository user.RepositoryUser
}

 func (a *userApp) CreateUser(ctx context.Context, nickname string, email string) (*user.User, error) {
	if nickname == "" || email == "" {
		return nil, ErrBadRequest
	}
	user := user.User{NickName: nickname, Email: email}

	userID, err := a.repository.AddUser(ctx, &user)

	if err != nil {
		return nil, err
	}

	user.UserID = userID

	return &user, nil
 }

 func (a *userApp) UpdateUser(ctx context.Context, nickname string, email string, userID int64, activate bool) (*user.User, error) {
	user, err := a.repository.UpdateUser(ctx, nickname, email, userID, activate)

	if err != nil {
		return nil, ErrNotFound
	}

	return user, nil
 }

 func (a *userApp) CheckUser(ctx context.Context, user_id int64) error {
	ok := a.repository.CheckUser(ctx, user_id)
	if !ok {
		return fmt.Errorf("not user in db")
	}
	return nil
 }

 func (a *userApp) GetUser(ctx context.Context, user_id int64) (*user.User, error) {
	user, err := a.repository.GetUser(ctx, user_id)

	if err != nil {
		return nil, ErrNotFound
	}
	
	return user, nil
 }

func (a *userApp) DeleteUser(ctx context.Context, user_id int64) (error) {
	err := a.repository.DeleteUser(ctx, user_id)

	if err != nil {
		return err
	}

	return nil
}

type UserDbApp interface {
	CreateUserDb(user user.UserDb) (*user.UserDb, error)
	GetUserDb(username, password string) (*user.UserDb, error)
	CheckUserDb(id int) (*user.UserDb, error)
	UpdateUserDb(username string, id int) (*user.UserDb, error)
	DeleteUserDb(id int) (*user.UserDb, error)
}

type authApp struct {
	repository user.RepositoryDbUser
}

func (a *authApp) CreateUserDb(user user.UserDb) (*user.UserDb, error) {
	user.Password = generatePasswordHash(user.Password)
	id, err := a.repository.CreateUserDb(user)
	user.Id = id
	return &user, err
}

const salt = "hjqrhjqw124617ajfhajs"

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (a *authApp) GetUserDb(username, password string) (*user.UserDb, error) {
	password = generatePasswordHash(password)
	u, err := a.repository.GetUser(username, password)
	return u, err
}

func (a *authApp) CheckUserDb(id int) (*user.UserDb, error) {
	return a.repository.CheckUserDb(id)
}

func (a *authApp) UpdateUserDb(username string, id int) (*user.UserDb, error) {
	return a.repository.UpdateUserDb(username, id)
}

func (a *authApp) DeleteUserDb(id int) (*user.UserDb, error) {
	return a.repository.DeleteUserDb(id)
}

func NewApp(repo ads.RepositryAd, repoUser user.RepositoryUser, repoUserDb user.RepositoryDbUser) App {
	return &appStruct{
		adApp: adApp{repository: repo},
		userApp: userApp{repository: repoUser},
		authApp: authApp{repository: repoUserDb},
	}
}
