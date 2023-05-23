package services

import (
	"context"
	"fmt"
	"strings"

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
	content := fmt.Sprintf("from IETF language code: %s\nto IETF language code: %s\ntranslate the text: '%s'\nproduce 3 versions of only the translation itself separated by a new line.", sourceLang, targetLang, text)
	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return strings.Split(strings.Replace(resp.Choices[0].Message.Content, "\"", "", -1), "\n"), err
}

func (s *OpenAITranslator) GenerateExamples(text, sourceLang string) ([]string, error) {
	content := fmt.Sprintf("in the language which has IETF language code: %s\ngive me examples of 3 sentences with the usage of '%s'. Separate sentences by one new line symbol, do not add quotes or sentence numbers.", sourceLang, text)
	resp, err := s.client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, line := range strings.Split(resp.Choices[0].Message.Content, "\n") {
		if line != "" {
			result = append(result, line)
		}
	}

	return result, err
}
