package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"povtorushka-bot/callbacks"
	"povtorushka-bot/services/telegram"
	"povtorushka-bot/types"
)

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

type APIGatewayRequest struct {
	Body string `json:"body"`
}

func Handler(ctx context.Context, request []byte) (*Response, error) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))

	vkConfirmationCode := os.Getenv("VK_CONFIRMATION_CODE")
	vkSecret := os.Getenv("VK_SECRET")
	vkGroupIDStr := os.Getenv("VK_GROUP_ID")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if vkConfirmationCode == "" {
		slog.Error("Environment variable VK_CONFIRMATION_CODE is required")
		os.Exit(1)
	}
	if vkSecret == "" {
		slog.Error("Environment variable VK_SECRET is required")
		os.Exit(1)
	}
	if vkGroupIDStr == "" {
		slog.Error("Environment variable VK_GROUP_ID is required")
		os.Exit(1)
	}
	vkGroupID, err := strconv.Atoi(vkGroupIDStr)
	if err != nil {
		slog.Error("Invalid VK_GROUP_ID", "error", err)
		os.Exit(1)
	}

	if telegramBotToken == "" {
		slog.Error("Environment variable TELEGRAM_BOT_TOKEN is required")
		os.Exit(1)
	}
	if telegramChatID == "" {
		slog.Error("Environment variable TELEGRAM_CHAT_ID is required")
		os.Exit(1)
	}

	var apiGatewayReq APIGatewayRequest

	if err := json.Unmarshal(request, &apiGatewayReq); err != nil {
		slog.Error("Error parsing API Gateway request", "error", err)
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid API Gateway request format",
		}, nil
	}

	vkRequestBody := []byte(apiGatewayReq.Body)

	slog.Debug("API Gateway request body", "body", string(vkRequestBody))

	var callback types.VKCallback
	if err := json.Unmarshal(vkRequestBody, &callback); err != nil {
		slog.Error("Error parsing VK callback JSON", "error", err, "body_preview", string(vkRequestBody[:min(100, len(vkRequestBody))]))
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body:       "Invalid JSON format",
		}, nil
	}

	slog.Debug("VK callback parsed", "callback", callback)

	if callback.Secret != vkSecret {
		slog.Error("Invalid secret key received",
			"received_length", len(callback.Secret),
			"expected_length", len(vkSecret))
		return &Response{
			StatusCode: http.StatusUnauthorized,
			Body:       "Invalid secret key",
		}, nil
	}

	if callback.Type == "confirmation" {
		if callback.GroupID != vkGroupID {
			slog.Error("Invalid group_id in confirmation request",
				"received", callback.GroupID,
				"expected", vkGroupID)
			return &Response{
				StatusCode: http.StatusBadRequest,
				Body:       "Invalid group_id",
			}, nil
		}

		slog.Info("Handling confirmation request", "group_id", callback.GroupID)
		return &Response{
			StatusCode: http.StatusOK,
			Body:       vkConfirmationCode,
		}, nil
	}

	tgService, err := telegram.New(telegramBotToken, telegramChatID)
	if err != nil {
		slog.Error("Error creating Telegram service", "error", err)
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       "Failed to initialize Telegram service",
		}, nil
	}
	switch callback.Type {

	case "wall_post_new":
		slog.Debug("Processing wall_post_new event", "object", callback.Object)
		if err := callbacks.HandleWallPostNew(callback.Object, tgService); err != nil {
			slog.Error("Error handling wall_post_new", "error", err)
		}

	default:
		slog.Warn("Unhandled event type", "type", callback.Type)
	}

	slog.Info("Return OK to VK")
	return &Response{
		StatusCode: http.StatusOK,
		Body:       "ok",
	}, nil
}

func main() {
}
