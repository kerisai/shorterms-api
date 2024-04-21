package config

import (
	"context"
	"errors"

	"github.com/google/generative-ai-go/genai"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
)

var (
	ErrEmptyGeminiApiKey = errors.New("empty gemini api key")
)

func CreateGeminiClient(geminiApiKey string) *genai.Client {
	if geminiApiKey == "" {
		log.Fatal().Err(ErrEmptyGeminiApiKey).Msg("failed to configure gemini client")
	}

	client, err := genai.NewClient(context.Background(), option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to configure gemini client")
	}

	return client
}
