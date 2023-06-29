package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"unicode"

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
		builder.WriteString(fmt.Sprintf("- %s", t.GetText()))
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

func IsExpressionWithAllData(e models.Expression) bool {
	return e.Translations != nil && e.Examples != nil && e.Audio != nil && len(e.Translations) > 0 && len(e.Examples) > 0 && len(e.Synonyms) > 0 && len(e.Audio) > 0
}

func BuildMessage(text ...string) string {
	var message string
	for i, t := range text {
		if t != "" {
			if i > 0 && message != "" {
				message += "\n\n"
			}
			message += t
		}
	}
	return message
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

func BuildOpenAiResponse(text string) []string {
	unique := make(map[string]bool)
	result := []string{}
	re := regexp.MustCompile(`["\d\.]+`)

	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		line = strings.TrimLeft(line, "-")
		line = re.ReplaceAllString(line, "")
		line = strings.TrimSpace(line)
		runes := []rune(line)
		runes[0] = unicode.ToUpper(runes[0])
		resultLine := string(runes)

		if !unique[resultLine] {
			result = append(result, resultLine)
			unique[resultLine] = true
		}

	}

	return result
}

func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func AddBlockTitleToText(title, text string) string {
	if text == "" || title == "" {
		return ""
	}
	return fmt.Sprintf("<b>%s</b>\n%s", title, text)
}
