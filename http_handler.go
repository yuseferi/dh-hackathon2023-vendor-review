package hackathon_2023_vendor_reviews

import (
	"context"
	"net/http"
	"time"

	"github.com/yuseferi/gocache"
	"hackathon-2023-vendor-reviews/rest"
)

var cache = gocache.NewCache(time.Second * 20)

func GetReviewsSummary(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	genAI := NewGenAI()
	summary, err := genAI.getOverAllSummary(ctx, getReviewsStatement(100))
	if err != nil {
		rest.Reply(w, http.StatusInternalServerError, err, nil)
	}

	rest.Reply(w, http.StatusOK, summary, nil)

}

func GetReviews(w http.ResponseWriter, r *http.Request) {
	const reviewNOLimit = 20
	const cacheKey = "reviews"
	ctx := context.Background()
	genAI := NewGenAI()
	reviews := make([]ReviewResponseItem, 0, 100)
	cachedReviews, found := cache.Get(cacheKey)
	if !found {
		suggestedReplies, err := genAI.getReviewsWithSmartReplies(ctx, getReviewsWithID(reviewNOLimit), getSmartRelies(20))
		//suggestedReplies, err := genAI.getReviewsWithSmartReplies(ctx, getReviewsWithID(10), getSmartRelies(5))
		if err != nil {
			rest.Reply(w, http.StatusInternalServerError, err, nil)
		}
		reviewsMap := getReviewsItems(reviewNOLimit)
		for _, item := range suggestedReplies {
			reviewsMap[item.ID].SuggestedReply = item.SuggestedReply
		}
		for _, item := range reviewsMap {
			reviews = append(reviews, *item)
		}
		cache.Set(cacheKey, reviews, time.Minute*5)
	} else {
		reviews = cachedReviews.([]ReviewResponseItem)
	}

	response := ReviewsResponse{
		HasPrevious: false,
		HasNext:     false,
		Reviews:     reviews,
		ReviewsDigest: struct {
			Tags []string `json:"tags"`
		}(struct{ Tags []string }{
			[]string{"test"},
		}),
	}
	//if err != nil {
	//	rest.Reply(w, http.StatusInternalServerError, err, nil)
	//}

	rest.Reply(w, http.StatusOK, response, nil)

}
