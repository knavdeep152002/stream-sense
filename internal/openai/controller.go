package openai

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/knavdeep152002/stream-sense/internal/utils"
	openai "github.com/sashabaranov/go-openai"
)

type query struct {
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
