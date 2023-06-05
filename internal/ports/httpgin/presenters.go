package httpgin

import (
	"time"

	"github.com/gin-gonic/gin"

	"ads/internal/ads"
	"ads/internal/user"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AuthorID  int64  `json:"author_id"`
	Published bool   `json:"published"`
	CreateDate time.Time `json:"create_date"`
	UpdateDate time.Time `json:"update_date"`
}

type userDeleteResponse struct {
	UserID int64 `json:"user_id"`
}

type deleteUserRequest struct {
	UserID int64 `json:"user_id"`
}

type getUserRequest struct {
	UserID int64 `json:"user_id"`
}

type adDeleteRequest struct {
	UserID int64 `json:"user_id"`
}

type adsRequest struct {
	Data []ads.Ad `json:"data"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type getAdRequest struct {
	ID int64 `json:"id"`
}

type createUserRequest struct {
	NickName string `json:"nickname"`
	Email   string `json:"email"`
}

type updateUserRequest struct {
	UserID  int64 `json:"user_id"`
	NickName string `json:"nickname"`
	Email   string `json:"email"`
	Activate bool	`json:"activate"`
}

type createUserResponse struct {
	UserID  int64 `json:"user_id"`
	NickName string `json:"nickname"`
	Email   string `json:"email"`
	Activate bool	`json:"activate"`
}

type createUserDB struct {
	Id       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
			CreateDate: ad.CreateDate,
			UpdateDate: ad.UpdateDate,
		},
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}

func AdsSuccessResponse(ads []*ads.Ad) *gin.H {
	result := []adResponse{}
	for _, ad := range ads {
		el := adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
			CreateDate: ad.CreateDate,
			UpdateDate: ad.UpdateDate,
		}
		result = append(result, el)
	}
	return &gin.H{
		"data": result,
	}
}

func UserSuccessResponse(u *user.User) *gin.H {
	return &gin.H{
		"data": createUserResponse{
			UserID: u.UserID,
			NickName: u.NickName,
			Email: u.Email,
			Activate: u.Activate,
		},
		"error": nil,
	}
}

func DeleteUserSuccess(user_id int64) *gin.H {
	return &gin.H{
		"data": userDeleteResponse{
			UserID: user_id,
		},
		"error": nil,
	}
}

func ErrUser(err error) *gin.H {
	return &gin.H{
		"error": err,
	}
}

func CreateUserDbSuccess(id int) *gin.H {
	return &gin.H{
		"id": id,
		"error": nil,
	}
}
