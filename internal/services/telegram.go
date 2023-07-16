package services

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	botTg "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/prplx/wordy/internal/types"
)

type TelegramService struct {
	bot  *tgbotapi.BotAPI
	bot2 *botTg.Bot
}

func NewTelegramService() *TelegramService {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot2, err := botTg.New(os.Getenv("TG_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	return &TelegramService{bot: bot, bot2: bot2}
}

func (s *TelegramService) SendText(chatId int64, text string, replyMessageId ...int) (string, error) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	replyTo := 0
	if len(replyMessageId) > 0 {
		replyTo = replyMessageId[0]
	}

	defer cancel()

	m, err := s.bot2.SendMessage(ctx, &botTg.SendMessageParams{
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
	var inlineKeyboardButtons [][]tgbotapi.InlineKeyboardButton
	for _, btn := range buttons {
		row := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.CallbackData))
		inlineKeyboardButtons = append(inlineKeyboardButtons, row)
	}

	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		inlineKeyboardButtons...,
	)
	response, err := s.bot.Send(msg)
	return response.Text, err
}

func (s *TelegramService) AnswerCallbackQuery(queryId string, text string) error {
	_, err := s.bot.AnswerCallbackQuery(tgbotapi.NewCallback(queryId, text))
	return err
}

func (s *TelegramService) SendTypingChatAction(chatId int64) error {
	_, err := s.bot.MakeRequest("sendChatAction", url.Values{
		"chat_id": {strconv.FormatInt(chatId, 10)},
		"action":  {"typing"},
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

	_, err := s.bot2.EditMessageText(ctx, params)
	return err
}

func (s *TelegramService) EditReplyMarkup(chatId int64, messageId int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	defer cancel()

	_, err := s.bot2.EditMessageReplyMarkup(ctx, &botTg.EditMessageReplyMarkupParams{
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

	_, err := s.bot2.DeleteMessage(ctx, &botTg.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: messageId,
	})
	return err
}
