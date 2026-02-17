package impl

import (
	"context"
	"log/slog"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/responses"
)

type ResponseDataBase struct {
	ctx    context.Context
	client *openai.Client
	model  openai.ResponsesModel
	logger *slog.Logger
}

type ResponseData struct {
	ResponseDataBase
	system param.Opt[string]
	input  responses.ResponseNewParamsInputUnion
	text   responses.ResponseTextConfigParam
}

func Response(data ResponseData) (string, error) {
	params := responses.ResponseNewParams{
		Model:        data.model,
		Input:        data.input,
		Instructions: data.system,
		Text:         data.text,
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
