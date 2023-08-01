package services

import (
	"context"
	"fmt"
	"time"

	"github.com/prplx/wordy/internal/helpers"
	openai "github.com/sashabaranov/go-openai"
)

type OpenAITranslator struct {
	client *openai.Client
}

func NewOpenAITranslator(apiKey string) *OpenAITranslator {
	return &OpenAITranslator{
		client: openai.NewClient(apiKey),
	}
}

func (s *OpenAITranslator) Translate(text, sourceLang, targetLang string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	content := fmt.Sprintf("From %s to %s translate text below: %q\nproduce 3 versions of only the translation itself separated by a new line. Do not add any notes or any extra information or context.", sourceLang, targetLang, text)
	resp, err := s.client.CreateChatCompletion(ctx, generateRequest(content))
	if err != nil {
		return nil, err
	}

	return helpers.BuildOpenAiResponse(resp.Choices[0].Message.Content), nil
}

func (s *OpenAITranslator) GenerateExamples(text, sourceLang string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	content := fmt.Sprintf("Generate 3 sentences in %s with examples of using the word %q. Separate sentences by one new line symbol, do not add quotes, dashes, or sentence numbers.", sourceLang, text)
	resp, err := s.client.CreateChatCompletion(ctx, generateRequest(content))
	if err != nil {
		return nil, err
	}

	return helpers.BuildOpenAiResponse(resp.Choices[0].Message.Content), nil
}

func (s *OpenAITranslator) GenerateSynonyms(text, sourceLang string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	content := fmt.Sprintf("Generate in %s 5 synonyms of the word %q. Separate synonyms by one new line symbol, do not add quotes, dashes, or synonym numbers.", sourceLang, text)
	resp, err := s.client.CreateChatCompletion(ctx, generateRequest(content))
	if err != nil {
		return nil, err
	}

	return helpers.BuildOpenAiResponse(resp.Choices[0].Message.Content), nil
}

func generateRequest(content string) openai.ChatCompletionRequest {
	return openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	}
}
