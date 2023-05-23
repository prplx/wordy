package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

func SetWebhookUrl(webhookUrl string) error {
	setWebhookUrl := "https://api.telegram.org/bot" + os.Getenv("TG_BOT_TOKEN") + "/setWebhook?url=" + webhookUrl + "/api/v1/bot"
	resp, err := http.Get(setWebhookUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response types.SetWebhookResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	if !response.Ok {
		return errors.New(response.Description)
	}

	fmt.Println(string(body))
	fmt.Println("Webhook URL is: " + webhookUrl + "/api/v1/bot")

	return nil
}

func IsProduction() bool {
	return Getenv("APP_ENV", "development") == "production"
}
