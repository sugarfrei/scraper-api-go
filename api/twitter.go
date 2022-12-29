package api

import (
	"net/http"
	"scraper-api-go/model"

	"github.com/gin-gonic/gin"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func (a *Api) Twitter(c *gin.Context) {

	scraper := twitterscraper.New()
	profile, err := scraper.GetProfile(c.Param("username"))
	if err != nil {
		a.abortWithError(c, model.ClientError("Bad user input").Prefix(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"Full Name": profile.Name, "Following": profile.FollowingCount, "Followers": profile.FollowersCount, "Tweets": profile.TweetsCount, "Likes": profile.LikesCount})
}
