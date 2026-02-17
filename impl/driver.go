package impl

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
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

	client := openai.NewClient(
		option.WithBaseURL(config.AiConfig.Url),
		option.WithAPIKey(config.AiConfig.Key),
	)
	logger.Info("OpenAI client initialized with URL", "url", config.AiConfig.Url)

	writerClient := aiClient{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(config.ModelsConfig.Writer),
		logger: logger,
	}

	classifierClient := aiClient{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(config.ModelsConfig.Classifier),
		logger: logger,
	}
	_ = writerClient

	skill := &skills[2]

	coverageResp, err := DetectCoverage(classifierClient, skill)
	if err != nil {
		return err
	}
	logger.Info("DetectCoverage response", "response", coverageResp)

	fitsResp, err := FitsSkill(classifierClient, skill, "I have experience with Go and Python.")
	if err != nil {
		return err
	}
	logger.Info("FitsSkill response", "response", fitsResp)

	return nil
}
