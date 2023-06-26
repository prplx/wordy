package services

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type TextToSpeechAWSPolly struct {
}

func NewTextToSpeechAWSPolly() *TextToSpeechAWSPolly {
	return &TextToSpeechAWSPolly{}
}

func (s *TextToSpeechAWSPolly) Convert(text, lang, userId string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{}))

	pollySvc := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(fmt.Sprintf("<speak><prosody rate='%d%%'>%s</prosody><break time='1s'/><prosody rate='%d%%'>%s</prosody></speak>", 90, text, 60, text)), TextType: aws.String("ssml"), VoiceId: aws.String("Joanna"), Engine: aws.String("neural")}

	output, err := pollySvc.SynthesizeSpeech(input)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s/%s.mp3", userId, time.Now().Format("2006-01-02-15-04-05"))
	s3params := &s3manager.UploadInput{
		Bucket:      aws.String("wordy-s3"),
		Key:         aws.String(fileName),
		Body:        output.AudioStream,
		ACL:         aws.String("public-read"),
		ContentType: aws.String("audio/mpeg"),
	}

	uploader := s3manager.NewUploader(sess)
	if _, err = uploader.Upload(s3params); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", os.Getenv("CDN_HOST_NAME"), fileName), nil
}
