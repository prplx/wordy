package services

import (
	"context"
	"fmt"
	"strings"
	"unicode"

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
	content := fmt.Sprintf("from IETF language code: %s\nto IETF language code: %s\ntranslate the text: '%s'\nproduce 3 versions of only the translation itself separated by a new line. Do not add any notes or any extra information or context.", sourceLang, targetLang, text)
	resp, err := s.client.CreateChatCompletion(context.Background(), generateRequest(content))
	if err != nil {
		return nil, err
	}

	return buildResponse(resp.Choices[0].Message.Content), nil
}

func (s *OpenAITranslator) GenerateExamples(text, sourceLang string) ([]string, error) {
	content := fmt.Sprintf("in the language which has IETF language code: %s\ngive me examples of 3 sentences with the usage of the word '%s'. Separate sentences by one new line symbol, do not add quotes, dashes, or sentence numbers.", sourceLang, text)
	resp, err := s.client.CreateChatCompletion(context.Background(), generateRequest(content))
	if err != nil {
		return nil, err
	}

	return buildResponse(resp.Choices[0].Message.Content), nil
}

func (s *OpenAITranslator) GenerateSynonyms(text, sourceLang string) ([]string, error) {
	content := fmt.Sprintf("in the language which has IETF language code: %s\ngive me 5 synonyms of '%s'. Separate synonyms by one new line symbol, do not add quotes, dashes, or synonym numbers.", sourceLang, text)
	resp, err := s.client.CreateChatCompletion(context.Background(), generateRequest(content))
	if err != nil {
		return nil, err
	}

	return buildResponse(resp.Choices[0].Message.Content), nil
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

func buildResponse(text string) []string {
	unique := make(map[string]bool)
	result := []string{}
	for _, line := range strings.Split(text, "\n") {
		if line != "" {
			trimmedLine := strings.TrimSpace(line)
			cleanLine := strings.TrimLeft(trimmedLine, "-")
			runes := []rune(cleanLine)
			runes[0] = unicode.ToUpper(runes[0])
			resultLine := string(runes)

			if !unique[resultLine] && line != "" {
				result = append(result, resultLine)
				unique[resultLine] = true
			}
		}
	}

	return result
}
