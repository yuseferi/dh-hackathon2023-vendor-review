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
	const reviewNOLimit = 50
	const cacheKey = "reviews"
	ctx := context.Background()
	genAI := NewGenAI()
	reviews := make([]ReviewResponseItem, 0, 100)
	var response ReviewsResponse
	cachedReviews, found := cache.Get(cacheKey)
	if !found {
		suggestedReplies, err := genAI.getReviewsWithSmartRepliesWithTags(ctx, getReviewsWithID(reviewNOLimit), getSmartRelies(reviewNOLimit/2))
		//suggestedReplies, err := genAI.getReviewsWithSmartRepliesWithTags(ctx, getReviewsWithID(10), getSmartRelies(5))
		if err != nil {
			rest.Reply(w, http.StatusInternalServerError, err, nil)
		}
		reviewsMap := getReviewsItems(reviewNOLimit)
		for _, item := range suggestedReplies {
			reviewsMap[item.ID].SuggestedReply = item.SuggestedReply
			reviewsMap[item.ID].Tags = item.Tags
		}
		for _, item := range reviewsMap {
			reviews = append(reviews, *item)
		}
		summary, err := genAI.getOverAllSummary(ctx, getReviewsStatement(reviewNOLimit))
		if err != nil {
			rest.Reply(w, http.StatusInternalServerError, err, nil)
		}
		response = ReviewsResponse{
			HasPrevious:   false,
			HasNext:       false,
			Reviews:       reviews,
			ReviewsDigest: summary,
		}
		cache.Set(cacheKey, response, time.Minute*50)
	} else {
		response = cachedReviews.(ReviewsResponse)
	}

	rest.Reply(w, http.StatusOK, response, nil)

}
