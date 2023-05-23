package services

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/prplx/wordy/internal/types"
)

type TelegramService struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramService() *TelegramService {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}
	return &TelegramService{bot: bot}
}

func (s *TelegramService) SendText(chatId int64, text string, replyMessageId int) (string, error) {
	msg := tgbotapi.NewMessage(chatId, text)
	if replyMessageId != 0 {
		msg.ReplyToMessageID = replyMessageId
	}
	response, err := s.bot.Send(msg)

	return response.Text, err
}

func (s *TelegramService) SendReplyKeyboard(chatId int64, buttons []types.KeyboardButton, text string) (string, error) {
	var inlineKeyboardButtons []tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		inlineKeyboardButtons = append(inlineKeyboardButtons, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.CallbackData))
	}

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			inlineKeyboardButtons...,
		),
	)
	response, err := s.bot.Send(msg)
	return response.Text, err
}

func (s *TelegramService) AnswerCallbackQuery(queryId string, text string) error {
	_, err := s.bot.AnswerCallbackQuery(tgbotapi.NewCallback(queryId, text))
	return err
}
