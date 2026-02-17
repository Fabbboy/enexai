package impl

import (
	"context"
	"log/slog"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ResponseData struct {
	ctx    context.Context
	client *openai.Client
	model  openai.ResponsesModel
	logger *slog.Logger
	input  responses.ResponseNewParamsInputUnion
}

func Response(data ResponseData) (string, error) {
	params := responses.ResponseNewParams{
		Model: data.model,
		Input: data.input,
	}

	start := time.Now()
	resp, err := data.client.Responses.New(
		data.ctx,
		params,
	)

	elapsed := time.Since(start)
	if err != nil {
		data.logger.Error("Error while getting response from OpenAI", "error", err)
		return "", err
	}

	duration := elapsed.Seconds()
	data.logger.Info("OpenAI response",
		"input_tokens", resp.Usage.InputTokens,
		"output_tokens", resp.Usage.OutputTokens,
		"duration_s", duration,
	)
	return resp.OutputText(), nil
}
