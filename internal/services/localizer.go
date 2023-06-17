package services

import "github.com/nicksnyder/go-i18n/v2/i18n"

type LocalizerService struct {
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
}

func NewLocalizerService(bundle *i18n.Bundle) *LocalizerService {
	return &LocalizerService{
		bundle:    bundle,
		localizer: i18n.NewLocalizer(bundle, "en"),
	}
}

func (s *LocalizerService) L(id string, count ...interface{}) string {
	var pluralCount interface{} = 1
	if len(count) > 0 {
		pluralCount = count[0]
	}

	return s.localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID:   id,
		PluralCount: pluralCount,
	})
}

func (s *LocalizerService) ChangeLanguage(lang string) {
	s.localizer = i18n.NewLocalizer(s.bundle, lang)
}
