package httpgin

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ads/internal/app"
)

func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			log.Println("error create ad", err)
			c.JSON(400, AdErrorResponse(err))
			return
		}

		if err := a.CheckUser(c, reqBody.UserID); err != nil {
			log.Println("not found user in db. Need create/register user")
			c.JSON(400, AdErrorResponse(err))
			return
		}

		ad, err := a.CreateAd(c.Request.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(500, AdErrorResponse(err))
			}
			c.JSON(200, AdErrorResponse(err))
			log.Println("error create ad")
			return
		}
		log.Println("Success create ad", http.StatusOK, "id ad", ad.ID)
		c.Status(http.StatusOK)
		c.JSON(200, AdSuccessResponse(ad))
		log.Default()
	}
}

func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			return
		}

		ad, err := a.ChangeAdStatus(c.Request.Context(), int64(adID), reqBody.Published, reqBody.UserID)
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(500, AdErrorResponse(err))
			}
			c.JSON(200, AdErrorResponse(err))
			log.Println("error change status err", err)
			return
		}
		log.Println("Success change status ad", http.StatusOK, "id ad", ad.ID)
		c.JSON(200, AdSuccessResponse(ad))
	}
}

func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error update ad", err)
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error update ad", err)
			return
		}

		ad, err := a.UpdateAd(c.Request.Context(), reqBody.UserID, reqBody.Title, reqBody.Text, int64(adID))
		if err != nil {
			if errors.Is(err, app.ErrForbidden) {
				c.JSON(403, AdErrorResponse(err))
			} else if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(500, AdErrorResponse(err))
			}
			c.JSON(200, AdErrorResponse(err))
			log.Println("error update err", err)
			return
		}
		log.Println("Success update ad", http.StatusOK, "id ad", ad.ID)
		c.JSON(200, AdSuccessResponse(ad))
	}
}

func getAd(a app.App, c *gin.Context, ad_id string) {
		var reqBody getAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error get ad", err)
			return
		}

		adID, err := strconv.Atoi(ad_id)
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error get ad", err)
			return
		}

		ad, err := a.GetAd(c.Request.Context(), int64(adID))
		if err != nil {
			if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} 
			c.Status(500)
			c.JSON(200, AdErrorResponse(err))
			log.Println("error get ad", err)
			return
		}
		log.Println("Success get ad", http.StatusOK, "id ad", ad.ID)
		c.JSON(200, AdSuccessResponse(ad))
}


func listAds(a app.App, c *gin.Context) {
		var reqBody adsRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error get ads", err)
			return
		}
		ads, err := a.ListAds(c)
		if err != nil {
			if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			}
			c.JSON(200, AdErrorResponse(err))
			log.Println("error get ads", err)
			return
		}
		log.Println("Success get ads", http.StatusOK)
		c.JSON(200, AdsSuccessResponse(ads))
}

func searchAdByName(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")

		ads, err := a.SearchAdByName(c, title)
		if err != nil {
			c.JSON(200, AdErrorResponse(err))
			log.Println("error get ads", err)
			return
		}
		log.Println("Success search ad", http.StatusOK, "id ad", ads[0].ID)
		c.JSON(200, AdsSuccessResponse(ads))
	}
}

func listAdsAuthor(a app.App, c *gin.Context) {
	authorID, err := strconv.Atoi(c.Query("author_id"))
	if err != nil {
		c.JSON(400, AdErrorResponse(err))
		log.Println("error get ads", err)
		return
	}
	ads, err := a.ListAdsAuthor(c, int64(authorID))
	if err != nil {
		if errors.Is(err, app.ErrBadRequest) {
			c.JSON(400, AdErrorResponse(err))
		}
		c.JSON(200, AdErrorResponse(err))
		log.Println("error get ads", err)
		return
	}
	log.Println("Success get ads filter: author", http.StatusOK, "author_id", ads[0].AuthorID)
	c.JSON(200, AdsSuccessResponse(ads))
}

func listAdsDate(a app.App, c *gin.Context) {
	d := c.Query("day")
	day, err := strconv.Atoi(d)
	if err != nil {
		c.JSON(400, AdErrorResponse(err))
		log.Println("error get ads", err)
		return
	}
	ads, err := a.ListAdsDate(c, int64(day))
	if err != nil {
		if errors.Is(err, app.ErrBadRequest) {
			c.JSON(400, AdErrorResponse(err))
		}
		c.JSON(200, AdErrorResponse(err))
		log.Println("error get ads", err)
		return
	}
	log.Println("Success get ads filter: day", http.StatusOK, "day", ads[0].CreateDate.Day())
	c.JSON(200, AdsSuccessResponse(ads))
}

func getAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		ad_id := c.Query("ad_id")
		if ad_id != "" {
			getAd(a, c, ad_id)
			return
		}
		filter := c.Query("filter")
		if filter == "author" {
			listAdsAuthor(a, c)
			return
		}
		if filter == "date" {
			listAdsDate(a, c)
			return
		}
		// default output ads
		listAds(a, c)
	}
}

func deleteAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody adDeleteRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error delete ad", err)
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error user update", err)
			return
		}

		ad, err := a.DeleteAd(c.Request.Context(), reqBody.UserID, int64(adID))
		if err != nil {
			c.Status(500)
			c.JSON(200, AdErrorResponse(err))
			log.Println("error delete ad", err)
			return
		}
		c.Status(http.StatusOK)
		c.JSON(200, AdSuccessResponse(ad))
		log.Default()
		log.Println("Success delete ad", http.StatusOK, "ad id", ad.ID, "user id", ad.AuthorID)
	}
}

func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createUserRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error create user", err)
			return
		}

		u, err := a.CreateUser(c.Request.Context(), reqBody.NickName, reqBody.Email)
		if err != nil {
			if errors.Is(err, app.ErrBadRequest) {
				c.JSON(400, AdErrorResponse(err))
			} else {
				c.JSON(500, AdErrorResponse(err))
			}
			c.JSON(200, AdErrorResponse(err))
			log.Println("error create user", err)
			return
		}
		c.Status(http.StatusOK)
		c.JSON(200, UserSuccessResponse(u))
		log.Default()
		log.Println("Success create user", http.StatusOK, "user id", u.UserID)
	}
}

func updateUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error user update", err)
			return
		}

		user_id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(400, AdErrorResponse(err))
			log.Println("error user update", err)
			return
		}

		u, err := a.UpdateUser(c.Request.Context(), reqBody.NickName, reqBody.Email, int64(user_id), reqBody.Activate)
		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(404, AdErrorResponse(err))
			}
			c.Status(500)
			c.JSON(200, AdErrorResponse(err))
			log.Println("error user update", err)
			return
		}
		log.Println("Success update user", http.StatusOK, "user id", u.UserID)
		c.JSON(200, UserSuccessResponse(u))
	}
}

func deleteUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody deleteUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, ErrUser(err))
			log.Println("error user delete", err)
			return
		}

		user_id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(400, ErrUser(err))
			log.Println("error user delete", err)
			return
		}

		if err := a.CheckUser(c, int64(user_id)); err != nil {
			log.Println("not found user in db. Need create/register user")
			c.JSON(400, AdErrorResponse(err))
			return
		}

		err = a.DeleteUser(c.Request.Context(), int64(user_id))
		if err != nil {
			c.Status(500)
			c.JSON(200, ErrUser(err))
			log.Println("error user delete", err)
			return
		}
		log.Println("Success delete user", http.StatusOK, "user id", user_id)
		c.JSON(200, DeleteUserSuccess(int64(user_id)))
	}
}

func getUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody getUserRequest
		if err := c.Bind(&reqBody); err != nil {
			c.JSON(400, ErrUser(err))
			log.Println("error user get", err)
			return
		}

		user_id, err := strconv.Atoi(c.Param("user_id"))
		if err != nil {
			c.JSON(400, ErrUser(err))
			log.Println("error user get", err)
			return
		}

		u, err := a.GetUser(c.Request.Context(), int64(user_id))
		if err != nil {
			if errors.Is(err, app.ErrNotFound) {
				c.JSON(404, AdErrorResponse(err))
			} 
			c.Status(500)
			c.JSON(200, ErrUser(err))
			log.Println("error user get", err)
			return
		}
		log.Println("Success get user", http.StatusOK, "user id", u.UserID)
		c.JSON(200, UserSuccessResponse(u))
	}
}
