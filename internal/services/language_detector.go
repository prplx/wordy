package services

type LanguageDetectorService struct {
	detector LanguageDetector
}

func NewLanguageDetectorService() *LanguageDetectorService {
	return &LanguageDetectorService{
		detector: NewLanguageDetectorLingua(),
	}
}

func (s *LanguageDetectorService) Detect(text string) (string, bool) {
	return s.detector.Detect(text)
}
