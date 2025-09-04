package callbacks

import (
	"encoding/json"
	"log/slog"
	"povtorushka-bot/services/telegram"
	"povtorushka-bot/types"
	"strconv"
)

func HandleWallPostNew(object json.RawMessage, tgService *telegram.Service) error {
	slog.Debug("Raw wall post object", "object", string(object))

	var post types.WallPost
	if err := json.Unmarshal(object, &post); err != nil {
		slog.Error("Error parsing wall post object", "error", err)
		return err
	}

	slog.Debug("Parsed wall post", "post", post)
	if post.PostType != "post" {
		slog.Warn("post_type != 'post'", "post_type", post.PostType)
		return nil
	}

	slog.Info("New wall post received",
		"id", post.ID,
		"owner_id", post.OwnerID,
		"from_id", post.FromID)

	var photoURLs []string
	slog.Debug("Processing attachments", "count", len(post.Attachments))
	for _, attachment := range post.Attachments {
		if attachment.Type == "photo" && attachment.Photo != nil && attachment.Photo.OrigPhoto != nil {
			photoURLs = append(photoURLs, attachment.Photo.OrigPhoto.URL)
		}
	}

	hasPhotos := len(photoURLs) > 0
	hasText := post.Text != ""

	if !hasPhotos && !hasText {
		slog.Info("Post has no content (text or photos), skipping Telegram notification")
		return nil
	}

	if hasText {
		slog.Debug("Post text", "length", len(post.Text), "text", post.Text)
	}
	slog.Info("Total photos to send", "count", len(photoURLs))

	// Формируем ссылку на автора
	var authorURL string
	if post.PostAuthorData != nil && post.PostAuthorData.Author != 0 {
		authorURL = "https://vk.com/id" + strconv.Itoa(post.PostAuthorData.Author)
		slog.Debug("Author URL generated", "author_id", post.PostAuthorData.Author, "url", authorURL)
	}

	if err := tgService.SendContent(post.Text, photoURLs, authorURL); err != nil {
		slog.Error("Error sending to Telegram", "error", err)
		return err
	}

	slog.Info("Message successfully sent to Telegram")
	return nil
}
