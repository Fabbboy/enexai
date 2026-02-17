package impl

import (
	"context"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type ResponseData struct {
	ctx    context.Context
	client *openai.Client
	model  openai.ResponsesModel
}

func Response(data ResponseData) (string, error) {
	params := responses.ResponseNewParams{
		Model: data.model,
		
	}

	resp, err := data.client.Responses.New(
		data.ctx,
		params,
	)

	if err != nil {
		return "", err
	}

	return resp.OutputText(), nil
}
