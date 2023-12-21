package hackathon_2023_vendor_reviews

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const Key = "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"

func NewGenAI() GenAI {
	return GenAI{
		client: openai.NewClient(Key),
	}
}

type GenAI struct {
	client *openai.Client
}
type openAISummaryResponse struct {
	Summary  string   `json:"summary"`
	Likes    []string `json:"likes"`
	Dislikes []string `json:"dislikes"`
}

type openAIReplyResponse struct {
	ID             string `json:"ID"`
	SuggestedReply string `json:"SuggestedReply"`
}

func (c GenAI) getOverAllSummary(ctx context.Context, reviews []string) (res openAISummaryResponse, err error) {
	systemPrompt := "Give the Response as a JSON, with Review Digest Summary and Positive Themes and Negative Themes separated. The Response must be in simple to understand English without complex compound words. use summary,likes,dislikes as keys for json reponse"
	summaryPrompt := "As a Restaurant Owner operating with an Online Food Delivery Platform, I received some reviews and ratings from my customers. Look at the Review Data Set and come up with a \"Reviews Digest Summary\" statement. The Summary statement must encompass the overall sentiments of the reviews to help the Restaurant owner gather insights and take actions, while also including some trend changes from previous to current period.\n\nA reference Reviews Digest Summary Statement is included below and the output must be not more than 2 sentences and 30 words.\nReviews Digest Summary Statement:\n\"12 customers added a review in the past month. 30% \n of them mentioned that food quality was not as expected. The number of reviews mentioning that items were missing from the order also increased by 15%\"\n\nAlso come up with the Top 3 Themes each around \"What Customers Love about you\" and \"What Customers don't Like about you\" where each Theme has to be 1 or 2 words. Some of the themes could be \"Food Taste\", \"Food Quality\", \"Food Portion\", \"Food Presentation\", \"Packaging\", \"Overall Experience\", \"Delivery Time\", \"Convenient\", \"Great Selection\", \"Sustainable Packaging\", \"Consistent\", \"Accommodating Requests\", \"Reliable\" etc."
	content := summaryPrompt + strings.Join(reviews, "\n")
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},

		Temperature: 0.2,
		//MaxTokens:   100,
		Stream: false,
	}
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()
	var responseStr string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		responseStr += response.Choices[0].Delta.Content
	}

	err = json.Unmarshal([]byte(responseStr), &res)
	if err != nil {
		return
	}
	return
}

func (c GenAI) getReviewsWithSmartReplies(ctx context.Context, reviews []string, sampleReplies []string) (res []openAIReplyResponse, err error) {
	systemPrompt := "based on given sample replies generate some good respectful replies for users reviews of a restaurant, response should be in a well-formed valid json format array"
	summaryPrompt := "first I add the reviews and sample replies. in the next message  I will provide a list of reviews that I need your reply suggestion for them, please return response in json format with keys, ID, suggestedReply"
	msg1 := summaryPrompt + "sample data: \n\n" + strings.Join(sampleReplies, "\n")
	msg2 := "and the data that I would like you generate suggested replies are  data: " + strings.Join(reviews, "\n")
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		//ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg1,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg2,
			},
		},

		Temperature: 0.2,
		//MaxTokens:   100,
		Stream: false,
	}
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()
	var responseStr string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		responseStr += response.Choices[0].Delta.Content
	}
	fmt.Println(len(responseStr))
	//responseStr = strings.ReplaceAll(responseStr, "\n", "")
	err = json.Unmarshal([]byte(responseStr), &res)
	if err != nil {
		return
	}
	return
}

func (c GenAI) getTags(ctx context.Context, reviews []string) (res []string, err error) {
	systemPrompt := "based on given sample replies generate some good respectful replies for users reviews of a restaurant, response should be in a well-formed valid json format array"
	summaryPrompt := "first I add the reviews and sample replies. in the next message  I will provide a list of reviews that I need your reply suggestion for them, please return response in json format with keys, ID, suggestedReply"
	msg1 := summaryPrompt + "sample data: \n\n" + strings.Join(reviews, "\n")
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		//ResponseFormat: &openai.ChatCompletionResponseFormat{Type: openai.ChatCompletionResponseFormatTypeJSONObject},
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: systemPrompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: msg1,
			},
		},

		Temperature: 0.2,
		//MaxTokens:   100,
		Stream: false,
	}
	stream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()
	var responseStr string
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("error:", err)
			break
		}
		responseStr += response.Choices[0].Delta.Content
	}
	err = json.Unmarshal([]byte(responseStr), &res)
	if err != nil {
		return
	}
	return
}
