package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/types"
)

func Getenv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func BuildMessageFromSliceOfTexted[WT types.WithText](wt []WT) string {
	var builder strings.Builder
	for i, t := range wt {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("%d. %s", i+1, t.GetText()))
	}

	return builder.String()
}

func GetUserFirstAndSecondLanguagesIds(user models.User, languages []models.Language) (models.Language, models.Language) {
	var firstLanguage models.Language
	var secondLanguage models.Language
	for _, language := range languages {
		if language.ID == uint(user.FirstLanguage) {
			firstLanguage = language
		}
		if language.ID == uint(user.SecondLanguage) {
			secondLanguage = language
		}
	}

	return firstLanguage, secondLanguage
}

func IsExpressionFulfilledWithTranslationsExamplesAudio(e models.Expression) bool {
	return e.ID != 0 && e.Translations != nil && e.Examples != nil && e.Audio != nil && len(e.Translations) > 0 && len(e.Examples) > 0 && len(e.Audio) > 0
}

func BuildMessage(text ...string) string {
	return strings.Join(text, "\n\n")
}
