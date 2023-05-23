package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type TextToSpeechPlayHT struct {
	userId    string
	secretKey string
}

type ConvertResponse struct {
	Status          string `json:"status"`
	TranscriptionId string `json:"transcriptionId"`
}

type ConvertStatusResponse struct {
	Converted    bool   `json:"converted"`
	AudioUrl     string `json:"audioUrl"`
	Error        bool   `json:"error"`
	ErrorMessage string `json:"errorMessage"`
}

func NewTextToSpeechPlayHT(userId, secretKey string) *TextToSpeechPlayHT {
	return &TextToSpeechPlayHT{
		userId:    userId,
		secretKey: secretKey,
	}
}

func (s *TextToSpeechPlayHT) Convert(text, lang string) (string, error) {
	response, err := s.convert(text, lang)
	if err != nil {
		return "", err
	}

	var conversionStatus ConvertStatusResponse
	for i := 0; i < 10; i++ {
		conversionStatus, err = s.getConversionStatus(response.TranscriptionId)
		if err != nil {
			return "", err
		}

		if conversionStatus.Converted {
			break
		}

		time.Sleep(1 * time.Second)
	}

	if !conversionStatus.Converted {
		return "", fmt.Errorf("conversion failed")
	}

	if conversionStatus.Error {
		return "", fmt.Errorf(conversionStatus.ErrorMessage)
	}

	return conversionStatus.AudioUrl, nil

}

func (s *TextToSpeechPlayHT) convert(text, lang string) (ConvertResponse, error) {
	var response ConvertResponse
	url := "https://play.ht/api/v1/convert"
	payload := fmt.Sprintf("{\"content\":[\"%s\"],\"voice\":\"%s\", \"globalSpeed\":\"%s\"}", text, getLang(lang), "90%")

	apiResponse, err := s.makeApiRequest(url, http.MethodPost, payload)
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(apiResponse, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *TextToSpeechPlayHT) getConversionStatus(transcriptionId string) (ConvertStatusResponse, error) {
	var response ConvertStatusResponse
	url := "https://play.ht/api/v1/articleStatus?transcriptionId=" + transcriptionId
	apiResponse, err := s.makeApiRequest(url, http.MethodGet, "")

	if err != nil {
		return response, err
	}

	err = json.Unmarshal(apiResponse, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (s *TextToSpeechPlayHT) makeApiRequest(url, method, payload string) ([]byte, error) {
	reader := strings.NewReader(payload)
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("AUTHORIZATION", s.secretKey)
	req.Header.Add("X-USER-ID", s.userId)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("error: %s", res.Status)
	}

	return body, nil
}

func getLang(lang string) string {
	switch lang {
	case "en":
		return "Matthew"
	case "ru":
		return "Maxim"
	case "nl":
		return "nl-NL-Standard-A"
	default:
		return "Matthew"
	}
}
