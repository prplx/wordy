package services

import "os"

type TranslatorService struct {
	client Translator
}

func NewTranslatorService() *TranslatorService {
	return &TranslatorService{
		client: NewOpenAITranslator(os.Getenv("OPENAI_API_KEY")),
	}
}

func (s *TranslatorService) Translate(text, sourceLang, targetLang string) ([]string, error) {
	return s.client.Translate(text, sourceLang, targetLang)
}

func (s *TranslatorService) GenerateExamples(text, sourceLang string) ([]string, error) {
	return s.client.GenerateExamples(text, sourceLang)
}
