package impl

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func Run(configPath string, skillsPath string) error {
	ctx := context.Background()
	tinthandler := tint.NewHandler(os.Stdout, nil)
	logger := slog.NewLogLogger(tinthandler, slog.LevelDebug)

	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	aiConfig := &config.aiConfig
	client := openai.NewClient(
		option.WithBaseURL(aiConfig.url),
		option.WithAPIKey(aiConfig.key),
		option.WithDebugLog(logger),
	)

	respData := ResponseData{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(aiConfig.model),
	}

	msg, err := Response(respData)
	if err != nil {
		return err
	}

	println(msg)
	return nil
}
