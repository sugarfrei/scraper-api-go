package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"scraper-api-go/model"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"k8s.io/klog"
)

const nextPageURL string = `https://www.instagram.com/graphql/query/?query_hash=%s&variables=%s`
const nextPagePayload string = `{"id":"%s","first":50,"after":"%s"}`

type Root struct {
	Data struct {
		User struct {
			ID             string `json:"id"`
			EdgeFollowedBy struct {
				Count int `json:"count"`
			} `json:"edge_followed_by"`
			EdgeFollow struct {
				Count int `json:"count"`
			} `json:"edge_follow"`
			FullName                 string `json:"full_name"`
			IsPrivate                bool   `json:"is_private"`
			Username                 string `json:"username"`
			EdgeOwnerToTimelineMedia struct {
				Count    int `json:"count"`
				PageInfo struct {
					HasNextPage bool   `json:"has_next_page"`
					EndCursor   string `json:"end_cursor"`
				} `json:"page_info"`
				Edges []struct {
					Node struct {
						EdgeMediaPreviewLike struct {
							Count int `json:"count"`
						} `json:"edge_media_preview_like,omitempty"`
						EdgeLikedBy struct {
							Count int `json:"count"`
						} `json:"edge_liked_by,omitempty"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"edge_owner_to_timeline_media"`
		}
	}
}

// func Parse(body []byte) (*Root, error) {
// 	var data Root
// 	err := json.Unmarshal(body, &data)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, edge := range data.Data.User.EdgeOwnerToTimelineMedia.Edges {
// 		if nextPage {
// 			likes += edge.Node.EdgeMediaPreviewLike.Count
// 		} else {
// 			likes += edge.Node.EdgeLikedBy.Count
// 		}
// 	}

// 	return &data, nil
// }

func (a *Api) Instagram(ctx *gin.Context) {
	ig := ctx.Param("username")

	var followers, following, likes, posts int
	var fullName, msg string
	var nextPage, private bool

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Instagram 219.0.0.12.117 Android")
		r.Headers.Set("X-Requested-With", "XMLHttpRequest")
		r.Headers.Set("Referer", "https://www.instagram.com/"+ig)
	})

	c.OnError(func(r *colly.Response, e error) {
		klog.Errorln("error:", e, r.Request.URL)
	})

	c.OnResponse(func(r *colly.Response) {

		var data Root
		err := json.Unmarshal(r.Body, &data)
		if err != nil {
			a.abortWithError(ctx, model.ServerError(err).Prefix("Problem while unmarshalling instagram data"))
			return
		}

		followers = data.Data.User.EdgeFollowedBy.Count
		following = data.Data.User.EdgeFollow.Count
		fullName = data.Data.User.FullName
		private = data.Data.User.IsPrivate
		posts = data.Data.User.EdgeOwnerToTimelineMedia.Count

		if !private {
			for _, edge := range data.Data.User.EdgeOwnerToTimelineMedia.Edges {
				if nextPage {
					likes += edge.Node.EdgeMediaPreviewLike.Count
				} else {
					likes += edge.Node.EdgeLikedBy.Count
				}
			}
		}

		if data.Data.User.EdgeOwnerToTimelineMedia.PageInfo.HasNextPage {
			nextPage = true
			nextPageVars := fmt.Sprintf(nextPagePayload, data.Data.User.ID, data.Data.User.EdgeOwnerToTimelineMedia.PageInfo.EndCursor)
			r.Request.Ctx.Put("variables", nextPageVars)
			u := fmt.Sprintf(
				nextPageURL,
				"69cba40317214236af40e7efa697781d",
				url.QueryEscape(nextPageVars),
			)

			r.Request.Visit(u)
		}

		if private {
			msg = fmt.Sprintf("Can not count %s's likes. Such a boring person, all mysterios...", data.Data.User.Username)
		} else {
			msg = fmt.Sprintf("Finally someone worth stalking, wish %s all the likes in this world!", data.Data.User.Username)
		}
	})

	c.Visit("https://www.instagram.com/api/v1/users/web_profile_info/?username=" + ig)

	ctx.JSON(http.StatusOK, gin.H{"Full Name": fullName, "Following": following, "Followers": followers, "Posts": posts, "Likes": likes, "Private": msg})
}
