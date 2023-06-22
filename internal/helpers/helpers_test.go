package helpers

import (
	"testing"

	"github.com/prplx/wordy/internal/models"
	"github.com/prplx/wordy/internal/types"
	"github.com/stretchr/testify/assert"
)

type testWithText struct {
	text string
}

func (t testWithText) GetText() string {
	return t.text
}

func TestBuildMessageFromSliceOfTexted(t *testing.T) {
	var tests = []struct {
		result   string
		expected string
	}{
		{
			BuildMessageFromSliceOfTexted([]types.WithText{
				testWithText{"test1"},
				testWithText{"test2"},
				testWithText{"test3"},
			}),
			"- test1\n- test2\n- test3",
		},
		{
			BuildMessageFromSliceOfTexted([]types.WithText{
				testWithText{"test1"},
			}),
			"- test1",
		},
		{
			BuildMessageFromSliceOfTexted([]types.WithText{}),
			"",
		},
	}

	for _, test := range tests {
		t.Run("BuildMessageFromSliceOfTexted", func(t *testing.T) {
			assert.Equal(t, test.expected, test.result)
		})
	}
}

func TestBuildMessage(t *testing.T) {
	var tests = []struct {
		result   string
		expected string
	}{
		{
			BuildMessage("test1", "test2", "test3"),
			"test1\n\ntest2\n\ntest3",
		},
		{
			BuildMessage("test1", "", "test3"),
			"test1\n\ntest3",
		},
		{
			BuildMessage("", "", ""),
			"",
		},
	}

	for _, test := range tests {
		t.Run("BuildMessage", func(t *testing.T) {
			assert.Equal(t, test.expected, test.result)
		})
	}
}

func TestIsExpressionWithAllData(t *testing.T) {
	var translations = []models.Translation{{Text: "hello"}, {Text: "hi"}}
	var examples = []models.Example{{Text: "Hello, world!"}, {Text: "Hi there!"}}
	var synonyms = []models.Synonym{{Text: "greetings"}, {Text: "salutations"}}
	var audio = []models.Audio{{Url: "http://example.com/audio.mp3"}}
	var tests = []struct {
		name     string
		result   bool
		expected bool
	}{
		{
			"All fields are present",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     examples,
				Synonyms:     synonyms,
				Audio:        audio,
			}),
			true,
		},
		{
			"Translations are missing",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     examples,
				Synonyms:     synonyms,
			}),
			false,
		},
		{
			"Examples are missing",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     examples,
				Audio:        audio,
			}),
			false,
		},
		{
			"Synonyms are missing",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Synonyms:     synonyms,
				Audio:        audio,
			}),
			false,
		},
		{
			"Audio is missing",
			IsExpressionWithAllData(models.Expression{
				Examples: examples,
				Synonyms: synonyms,
				Audio:    audio,
			}),
			false,
		},
		{
			"Translations are empty",
			IsExpressionWithAllData(models.Expression{
				Translations: []models.Translation{},
				Examples:     examples,
				Synonyms:     synonyms,
				Audio:        audio,
			}),
			false,
		},
		{
			"Examples are empty",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     []models.Example{},
				Synonyms:     synonyms,
				Audio:        audio,
			}),
			false,
		},
		{
			"Synonyms are empty",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     examples,
				Synonyms:     []models.Synonym{},
				Audio:        audio,
			}),
			false,
		},
		{
			"Audio is empty",
			IsExpressionWithAllData(models.Expression{
				Translations: translations,
				Examples:     examples,
				Synonyms:     synonyms,
				Audio:        []models.Audio{},
			}),
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.result)
		})
	}
}
