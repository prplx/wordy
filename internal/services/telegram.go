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

func (s *TelegramService) SendText(chatId int64, text string, replyMessageId ...int) (string, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	replyTo := 0
	if len(replyMessageId) > 0 {
		replyTo = replyMessageId[0]
	}

	defer cancel()

	m, err := s.bot.SendMessage(ctx, &botTg.SendMessageParams{
		ChatID:           chatId,
		Text:             text,
		ParseMode:        models.ParseModeHTML,
		ReplyToMessageID: replyTo,
	})

	if err != nil {
		return "", err
	}

	return m.Text, nil
}

func (s *TelegramService) SendReplyKeyboard(chatId int64, buttons []types.KeyboardButton, text string) (string, error) {
	var inlineKeyboardButtons [][]models.InlineKeyboardButton
	for _, btn := range buttons {
		row := []models.InlineKeyboardButton{{Text: btn.Text, CallbackData: btn.CallbackData}}
		inlineKeyboardButtons = append(inlineKeyboardButtons, row)
	}

	msg := &botTg.SendMessageParams{
		ChatID:    chatId,
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

func (s *TelegramService) SendTypingChatAction(chatId int64) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()
	_, err := s.bot.SendChatAction(ctx, &botTg.SendChatActionParams{
		ChatID: chatId,
		Action: "typing",
	})
	return err
}

func (s *TelegramService) EditMessage(chatId int64, messageId int, text string, buttons ...[]types.KeyboardButton) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	params := &botTg.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: messageId,
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

func (s *TelegramService) EditReplyMarkup(chatId int64, messageId int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	_, err := s.bot.EditMessageReplyMarkup(ctx, &botTg.EditMessageReplyMarkupParams{
		ChatID:    chatId,
		MessageID: messageId,
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{models.InlineKeyboardButton{Text: "Hello", CallbackData: "hello"}}},
		},
	})
	return err
}

func (s *TelegramService) DeleteMessage(chatId int64, messageId int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	_, err := s.bot.DeleteMessage(ctx, &botTg.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: messageId,
	})
	return err
}
