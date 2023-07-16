package services

import (
	"github.com/pemistahl/lingua-go"
	"github.com/prplx/wordy/internal/helpers"
)

type LanguageDetectorLingua struct {
	detector lingua.LanguageDetector
}

func NewLanguageDetectorLingua() *LanguageDetectorLingua {
	var linguaLangs []lingua.Language

	for _, lang := range helpers.GetLanguageMap() {
		linguaLangs = append(linguaLangs, lang.LinguaLang)
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
