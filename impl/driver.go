package impl

import (
	"bufio"
	"context"
	"fmt"
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
	if err := config.Validate(); err != nil {
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

func Run(configPath, skillsPath string, debug bool) error {
	ctx := context.Background()
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level: level,
	})
	logger := slog.New(handler)

	config, skills, err := loadData(logger, configPath, skillsPath)
	if err != nil {
		return err
	}

	client := openai.NewClient(
		option.WithBaseURL(config.ApiConfig.Url),
		option.WithAPIKey(config.ApiConfig.Key),
	)
	logger.Info("OpenAI client initialized with URL", "url", config.ApiConfig.Url)

	classifierClient := aiClient{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(config.ModelsConfig.Classifier),
		logger: logger,
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Feedback title: ")
	scanner.Scan()
	title := scanner.Text()

	fmt.Print("Evidence: ")
	scanner.Scan()
	evidence := scanner.Text()

	text := title + "\n" + evidence
	logger.Info("Finding fitting skills", "title", title)

	matches, err := FindFittingSkills(classifierClient, skills, text)
	if err != nil {
		return err
	}

	fmt.Printf("\n%d skill(s) matched:\n", len(matches))
	for _, m := range matches {
		fmt.Printf("  - [%s] %s\n", m.Fitness, skills[m.Index].Competence)
	}

	return nil
}
