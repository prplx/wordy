package services

import (
	"github.com/pemistahl/lingua-go"
)

type LanguageDetectorLingua struct {
	detector lingua.LanguageDetector
}

func NewLanguageDetectorLingua() *LanguageDetectorLingua {
	linguaLangs := []lingua.Language{
		lingua.English,
		lingua.Russian,
		lingua.Dutch,
	}

	return &LanguageDetectorLingua{
		detector: lingua.NewLanguageDetectorBuilder().
			FromLanguages(linguaLangs...).
			Build(),
	}
}

func (s *LanguageDetectorLingua) Detect(text string) (string, bool) {
	lang, exists := s.detector.DetectLanguageOf(text)
	return lang.String(), exists
}
