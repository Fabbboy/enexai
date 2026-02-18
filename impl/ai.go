package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"text/template"
	"time"

	"github.com/openai/openai-go/v3"
)

type aiClient struct {
	ctx    context.Context
	client *openai.Client
	model  openai.ChatModel
	logger *slog.Logger
}

func (c aiClient) Send(params openai.ChatCompletionNewParams) (string, error) {
	params.Model = c.model

	start := time.Now()
	resp, err := c.client.Chat.Completions.New(
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
		"input_tokens", resp.Usage.PromptTokens,
		"output_tokens", resp.Usage.CompletionTokens,
		"duration_s", duration,
	)

	if len(resp.Choices) == 0 {
		return "", nil
	}
	return resp.Choices[0].Message.Content, nil
}

func jsonSchemaFormat(name string, schema map[string]any) openai.ChatCompletionNewParamsResponseFormatUnion {
	return openai.ChatCompletionNewParamsResponseFormatUnion{
		OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
			JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
				Name:   name,
				Schema: schema,
				Strict: openai.Bool(true),
			},
		},
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
	funcs := template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}
	t, err := template.New("").Funcs(funcs).Parse(tmpl)
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
