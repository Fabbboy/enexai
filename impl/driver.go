package impl

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/responses"
)

func loadData(logger *slog.Logger, configPath, skillsPath string) (*Config, []Skill, error) {
	logger.Info("Loading config from", "path", configPath)
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Config loaded")

	logger.Info("Loading skills from", "path", skillsPath)
	skills, err := LoadSkills(skillsPath)
	if err != nil {
		return nil, nil, err
	}
	logger.Info("Skills loaded", "count", len(skills))

	return config, skills, nil
}

func Run(configPath, skillsPath string) error {
	ctx := context.Background()
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	config, skills, err := loadData(logger, configPath, skillsPath)
	if err != nil {
		return err
	}
	_ = skills

	client := openai.NewClient(
		option.WithBaseURL(config.AiConfig.Url),
		option.WithAPIKey(config.AiConfig.Key),
	)
	logger.Info("OpenAI client initialized with URL", "url", config.AiConfig.Url)

	question := "What are the top 5 skills for a software engineer?"
	input := responses.ResponseNewParamsInputUnion{
		OfString: openai.String(question),
	}

	respData := ResponseData{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(config.AiConfig.Model),
		logger: logger,
		input:  input,
	}

	msg, err := Response(respData)
	if err != nil {
		return err
	}

	logger.Info("Response received", "message", msg)
	return nil
}
