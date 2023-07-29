package services

import (
	"context"
	"os"
	"os/signal"

	botTg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/prplx/wordy/internal/types"
)

type TelegramService struct {
	bot *botTg.Bot
}

func NewTelegramService() *TelegramService {
	bot, err := botTg.New(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	return &TelegramService{bot: bot}
}

func (s *TelegramService) SendText(chatID int64, text string, replyMessageID ...int) (string, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	replyTo := 0
	if len(replyMessageID) > 0 {
		replyTo = replyMessageID[0]
	}

	defer cancel()

	m, err := s.bot.SendMessage(ctx, &botTg.SendMessageParams{
		ChatID:           chatID,
		Text:             text,
		ParseMode:        models.ParseModeHTML,
		ReplyToMessageID: replyTo,
	})

	if err != nil {
		return "", err
	}

	return m.Text, nil
}

func (s *TelegramService) SendReplyKeyboard(chatID int64, buttons []types.KeyboardButton, text string) (string, error) {
	var inlineKeyboardButtons [][]models.InlineKeyboardButton
	for _, btn := range buttons {
		row := []models.InlineKeyboardButton{{Text: btn.Text, CallbackData: btn.CallbackData}}
		inlineKeyboardButtons = append(inlineKeyboardButtons, row)
	}

	msg := &botTg.SendMessageParams{
		ChatID:    chatID,
		Text:      text,
		ParseMode: models.ParseModeHTML,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboardButtons,
		},
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	m, err := s.bot.SendMessage(ctx, msg)
	if err != nil {
		return "", err
	}

	return m.Text, nil
}

func (s *TelegramService) SendTypingChatAction(chatID int64) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()
	_, err := s.bot.SendChatAction(ctx, &botTg.SendChatActionParams{
		ChatID: chatID,
		Action: "typing",
	})
	return err
}

func (s *TelegramService) EditMessage(chatID int64, messageID int, text string, buttons ...[]types.KeyboardButton) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	params := &botTg.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: messageID,
		Text:      text,
	}

	var inlineKeyboardButtons [][]models.InlineKeyboardButton
	if len(buttons) > 0 {
		for _, btn := range buttons[0] {
			row := []models.InlineKeyboardButton{{Text: btn.Text, CallbackData: btn.CallbackData}}
			inlineKeyboardButtons = append(inlineKeyboardButtons, row)
		}
	}
	if len(inlineKeyboardButtons) > 0 {
		params.ReplyMarkup = models.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboardButtons,
		}
	}

	_, err := s.bot.EditMessageText(ctx, params)
	return err
}

func (s *TelegramService) EditReplyMarkup(chatID int64, messageID int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	_, err := s.bot.EditMessageReplyMarkup(ctx, &botTg.EditMessageReplyMarkupParams{
		ChatID:    chatID,
		MessageID: messageID,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{models.InlineKeyboardButton{Text: "Hello", CallbackData: "hello"}}},
		},
	})
	return err
}

func (s *TelegramService) DeleteMessage(chatID int64, messageID int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	_, err := s.bot.DeleteMessage(ctx, &botTg.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: messageID,
	})
	return err
}
