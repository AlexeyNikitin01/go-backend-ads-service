package httpgin

import (
	"github.com/gin-gonic/gin"

	"ads/internal/app"
)

func AppRouter(r *gin.RouterGroup, a app.App) {
	r.GET("/ads", getAds(a))
	r.GET("/ads/search", searchAdByName(a))
	r.PUT("/ads/:ad_id/status", changeAdStatus(a))
	r.PUT("/ads/:ad_id", updateAd(a))
	r.POST("/ads", createAd(a))
	r.DELETE("/ads/delete/:ad_id", deleteAd(a))

	r.POST("/user", createUser(a))
	r.PUT("/user/update/:user_id", updateUser(a))
	r.DELETE("/user/delete/:user_id", deleteUser(a))
	r.GET("/user/:user_id", getUser(a))

	r.POST("/sign-up", signUp(a))
	r.POST("/sign-in", signIn(a))
	r.POST("/check-user", checkUserDb(a))
	r.POST("/delete-user", deleteUserDb(a))
	r.POST("/update-user", updateUserDb(a))
}
