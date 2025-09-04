package telegram

import (
	"log/slog"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Service struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func New(botToken, chatID string) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		return nil, err
	}

	bot.Debug = os.Getenv("DEBUG") == "true"
	slog.Info("Telegram bot authorized", "username", bot.Self.UserName)

	chatIDInt, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		slog.Error("Failed to parse chat ID as number", "error", err)
		return nil, err
	}

	return &Service{
		bot:    bot,
		chatID: chatIDInt,
	}, nil
}

func (s *Service) SendText(text string) error {
	slog.Info("Sending text message to Telegram")
	slog.Debug("Text message details", "chat_id", s.chatID, "text", text)

	msg := tgbotapi.NewMessage(s.chatID, text)
	msg.ParseMode = "Markdown"

	sentMsg, err := s.bot.Send(msg)
	if err != nil {
		return err
	}

	slog.Info("Text message sent successfully", "message_id", sentMsg.MessageID)
	return nil
}

func (s *Service) SendPhotosWithText(text string, photoURLs []string) error {
	slog.Info("Sending photos as media group to Telegram", "photo_count", len(photoURLs))
	slog.Debug("Media group details", "chat_id", s.chatID, "text", text, "photo_urls", photoURLs)

	var mediaGroup []interface{}

	for i, photoURL := range photoURLs {
		inputMedia := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(photoURL))
		if i == 0 && text != "" {
			inputMedia.Caption = text
			inputMedia.ParseMode = "Markdown"
		}

		mediaGroup = append(mediaGroup, inputMedia)
	}

	msg := tgbotapi.NewMediaGroup(s.chatID, mediaGroup)
	slog.Debug("Prepared media group", "media_count", len(mediaGroup))

	messages, err := s.bot.SendMediaGroup(msg)
	if err != nil {
		slog.Error("Failed to send media group", "error", err)
		return err
	}

	slog.Info("Media group sent successfully", "messages_sent", len(messages))
	return nil
}

func (s *Service) SendContent(text string, photoURLs []string, authorURL string) error {
	hasPhotos := len(photoURLs) > 0
	hasText := text != ""

	slog.Debug("SendContent called", "has_text", hasText, "has_photos", hasPhotos, "photo_count", len(photoURLs), "author_url", authorURL)

	if !hasPhotos && !hasText {
		slog.Info("No content to send (neither text nor photos)")
		return nil
	}

	if authorURL != "" {
		text += "\n\n[Связаться с автором](" + authorURL + ")"
	}

	if hasPhotos {
		return s.SendPhotosWithText(text, photoURLs)
	}

	return s.SendText(text)
}
