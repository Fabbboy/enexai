package impl

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

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

	if len(matches) == 0 {
		fmt.Println("\nNo skills matched.")
		return nil
	}

	fmt.Printf("\n%d skill(s) matched:\n", len(matches))
	for i, m := range matches {
		fmt.Printf("  %d) [%s] %s\n", i+1, m.Fitness, skills[m.Index].Competence)
	}

	fmt.Print("\nSelect skills (comma-separated numbers, e.g. 1,3): ")
	scanner.Scan()
	input := scanner.Text()

	var selected []SkillMatch
	for _, part := range strings.Split(input, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || n < 1 || n > len(matches) {
			return fmt.Errorf("invalid selection: %s", part)
		}
		selected = append(selected, matches[n-1])
	}

	writerClient := aiClient{
		ctx:    ctx,
		client: &client,
		model:  openai.ResponsesModel(config.ModelsConfig.Writer),
		logger: logger,
	}

	for _, m := range selected {
		skill := &skills[m.Index]
		fmt.Printf("\n--- %s ---\n", skill.Competence)

		logger.Info("Analyzing style", "competence", skill.Competence)
		style, err := AnalyzeStyle(classifierClient, skill)
		if err != nil {
			return err
		}

		logger.Info("Detecting coverage", "competence", skill.Competence)
		coverage, err := DetectCoverage(classifierClient, skill)
		if err != nil {
			return err
		}

		logger.Info("Writing evidence", "competence", skill.Competence)
		evidence, err := WriteEvidence(writerClient, skill, title, text, style, coverage)
		if err != nil {
			return err
		}

		fmt.Printf("\n%s\n", evidence)
	}

	return nil
}
