package services

type TextToSpeechService struct {
	client TextToSpeech
}

func NewTextToSpeechService() *TextToSpeechService {
	return &TextToSpeechService{
		// client: NewTextToSpeechPlayHT(os.Getenv("PLAYHT_USER_ID"), os.Getenv("PLAYHT_SECRET_KEY")),
		client: NewTextToSpeechAWSPolly(),
	}
}

func (s *TextToSpeechService) Convert(text, lang, userID string) (string, error) {
	return s.client.Convert(text, lang, userID)
}
