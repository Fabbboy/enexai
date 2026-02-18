package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"text/template"
	"time"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

type aiClient struct {
	ctx    context.Context
	client *openai.Client
	model  openai.ResponsesModel
	logger *slog.Logger
}

func (c aiClient) Send(params responses.ResponseNewParams) (string, error) {
	params.Model = c.model

	start := time.Now()
	resp, err := c.client.Responses.New(
		c.ctx,
		params,
	)

	elapsed := time.Since(start)
	if err != nil {
		c.logger.Error("Error while getting response from OpenAI", "error", err)
		return "", err
	}

	duration := elapsed.Seconds()
	c.logger.Info("OpenAI response",
		"input_tokens", resp.Usage.InputTokens,
		"output_tokens", resp.Usage.OutputTokens,
		"duration_s", duration,
	)
	return resp.OutputText(), nil
}

func textMsg(text string) responses.ResponseInputItemUnionParam {
	return responses.ResponseInputItemParamOfMessage(text, responses.EasyInputMessageRoleUser)
}

func inputItems(items ...responses.ResponseInputItemUnionParam) responses.ResponseNewParamsInputUnion {
	return responses.ResponseNewParamsInputUnion{OfInputItemList: responses.ResponseInputParam(items)}
}

func jsonSchemaFormat(name string, schema map[string]any) responses.ResponseTextConfigParam {
	return responses.ResponseTextConfigParam{
		Format: responses.ResponseFormatTextConfigParamOfJSONSchema(name, schema),
	}
}

func parse[T any](raw string) (T, error) {
	var result T
	err := json.Unmarshal([]byte(raw), &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

func renderTemplate(tmpl string, data any) (string, error) {
	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
