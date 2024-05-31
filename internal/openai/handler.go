package openai

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/knavdeep152002/stream-sense/internal/utils"
	openai "github.com/sashabaranov/go-openai"
)

type Query struct {
	Question string `json:"question"`
}

func qa(videoId string, question string) (content string, err error, statusCode int) {
	client := openai.NewClient(utils.OPENAI_API_KEY)
	// Read transcript from file
	fileContent, err := os.ReadFile(path.Join(utils.ROOT_TRANSCRIPT_DIR, fmt.Sprintf("transcribed_text-%s.txt", videoId)))
	if err != nil {
		statusCode = 404
		fmt.Println("Error reading file", err)
		return
	}
	systemPrompt := fmt.Sprintf(
		`You are a Context Helper AI, You answer questions based on the context provided.
		Here is the context: %s.
		Now, answer the following question: %s`, fileContent, question)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: systemPrompt,
				},
			},
		},
	)
	if err != nil {
		fmt.Println("Error completing chat", err)
		statusCode = 503
		return
	}
	if len(resp.Choices) == 0 {
		statusCode = 500
		err = fmt.Errorf("No choices found")
		return
	}
	content = resp.Choices[0].Message.Content
	return
}

// @Summary		Answer a question based on the video context
// @Description	Answer a question based on the video context
// @Tags			OpenAI
// @Accept			json
// @Produce		json
// @Param			videoId		path		string	true	"Video ID"
// @Param			question	body		Query	true	"Question"
// @Success		200			{object}	string
// @Failure		400			{object}	string
// @Failure		404			{object}	string
// @Failure		500			{object}	string
// @Router			/qa/{videoId} [POST]
func VideoIntelligence(c *gin.Context) {
	vidID := c.Param("videoId")
	question := &Query{}
	if err := c.BindJSON(question); err != nil {
		c.JSON(400, gin.H{"message": "error in binding json : " + err.Error()})
		return
	}
	log.Println("Question: ", question.Question)
	log.Println("Video ID: ", vidID)
	ans, err, status := qa(vidID, question.Question)
	if err != nil {
		c.JSON(status, gin.H{"message": "error in answering question : " + err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": ans})
}
