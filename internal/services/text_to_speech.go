package services

import "os"

type TextToSpeechService struct {
	client TextToSpeech
}

func NewTextToSpeechService() *TextToSpeechService {
	return &TextToSpeechService{
		client: NewTextToSpeechPlayHT(os.Getenv("PLAYHT_USER_ID"), os.Getenv("PLAYHT_SECRET_KEY")),
	}
}

func (s *TextToSpeechService) Convert(text, lang string) (string, error) {
	return s.client.Convert(text, lang)
}
