package services

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/rs/xid"
)

type TextToSpeechAWSPolly struct {
}

func NewTextToSpeechAWSPolly() *TextToSpeechAWSPolly {
	return &TextToSpeechAWSPolly{}
}

func (s *TextToSpeechAWSPolly) Convert(text, lang string) (string, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{}))

	pollySvc := polly.New(sess)

	input := &polly.SynthesizeSpeechInput{OutputFormat: aws.String("mp3"), Text: aws.String(fmt.Sprintf("<speak><prosody rate='%d%%'>%s</prosody><break time='2s'/><prosody rate='%d%%'>%s</prosody></speak>", 90, text, 50, text)), TextType: aws.String("ssml"), VoiceId: aws.String("Joanna"), Engine: aws.String("neural")}

	output, err := pollySvc.SynthesizeSpeech(input)
	if err != nil {
		return "", err
	}

	guid := xid.New()
	fileName := guid.String() + ".mp3"
	s3params := &s3manager.UploadInput{
		Bucket: aws.String("wordy-s3"),
		Key:    aws.String(fileName),
		Body:   output.AudioStream,
		ACL:    aws.String("public-read"),
	}

	uploader := s3manager.NewUploader(sess)
	if _, err = uploader.Upload(s3params); err != nil {
		return "", err
	}

	return fmt.Sprintf("https://wordy-s3.s3.eu-central-1.amazonaws.com/%s", fileName), nil
}
