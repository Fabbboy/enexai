package impl

import (
	_ "embed"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/openai/openai-go/v3"
)

//go:embed prompts/fits_skill.txt
var fitsSkillPrompt string

//go:embed prompts/style_analysis.txt
var styleAnalysisPrompt string

//go:embed prompts/evidence_analysis.txt
var evidenceAnalysisPrompt string

//go:embed prompts/coverage_detection.txt
var coverageDetectionPrompt string

//go:embed prompts/write_evidence.txt
var writeEvidencePrompt string

type Skill struct {
	Category   string `csv:"Kompetenzkategorie"`
	Competence string `csv:"Kompetenz"`
	Note       string `csv:"Bemerkung"`
}

func (s *Skill) FormatContext() string {
	var b strings.Builder
	b.WriteString("Category: ")
	b.WriteString(s.Category)
	b.WriteString("\nCompetence: ")
	b.WriteString(s.Competence)
	b.WriteString("\nNote: ")
	b.WriteString(s.Note)
	return b.String()
}

func (s *Skill) evidenceText() string {
	text := strings.TrimSpace(s.Note)
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	lines := strings.Split(text, "\n")

	var b strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			if b.Len() > 0 {
				b.WriteString(" ")
			}
			b.WriteString(line)
		}
	}

	return b.String()
}

func LoadSkills(file string) ([]Skill, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var skills []Skill
	err = gocsv.UnmarshalFile(f, &skills)
	if err != nil {
		return nil, err
	}

	return skills, nil
}

type Fitness string

const (
	Fit     Fitness = "fit"
	WeakFit Fitness = "weak_fit"
	NoFit   Fitness = "no_fit"
)

type FitsResult struct {
	Fitness Fitness `json:"fit"`
	Reason  string  `json:"reason"`
}

var fitsResultSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"fit": map[string]any{
			"type": "string",
			"enum": []string{string(Fit), string(WeakFit), string(NoFit)},
		},
		"reason": map[string]any{
			"type": "string",
		},
	},
	"required":             []string{"fit", "reason"},
	"additionalProperties": false,
}

func FitsSkill(client aiClient, skill *Skill, text string) (*FitsResult, error) {
	instructions, err := renderTemplate(fitsSkillPrompt, skill)
	if err != nil {
		return nil, err
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(text),
		},
		ResponseFormat: jsonSchemaFormat("fits_result", fitsResultSchema),
	}

	resp, err := client.Send(params)
	if err != nil {
		return nil, err
	}

	result, err := parse[FitsResult](resp)
	if err != nil {
		return nil, err
	}

	client.logger.Debug("FitsSkill", "competence", skill.Competence, "fitness", result.Fitness, "reason", result.Reason)
	return &result, nil
}

type SkillMatch struct {
	Index   int
	Fitness Fitness
}

func FindFittingSkills(client aiClient, skills []Skill, text string) ([]SkillMatch, error) {
	var matches []SkillMatch
	for i := range skills {
		result, err := FitsSkill(client, &skills[i], text)
		if err != nil {
			return nil, err
		}
		if result.Fitness == Fit || result.Fitness == WeakFit {
			matches = append(matches, SkillMatch{Index: i, Fitness: result.Fitness})
		}
	}
	return matches, nil
}

type StyleResult struct {
	Language       string `json:"language"`
	Perspective    string `json:"perspective"`
	Tense          string `json:"tense"`
	Tone           string `json:"tone"`
	SentenceLength string `json:"sentence_length"`
	UsesReferences bool   `json:"uses_references"`
}

var styleResultSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"language":        map[string]any{"type": "string"},
		"perspective":     map[string]any{"type": "string"},
		"tense":           map[string]any{"type": "string"},
		"tone":            map[string]any{"type": "string"},
		"sentence_length": map[string]any{"type": "string"},
		"uses_references": map[string]any{"type": "boolean"},
	},
	"required":             []string{"language", "perspective", "tense", "tone", "sentence_length", "uses_references"},
	"additionalProperties": false,
}

func AnalyzeStyle(client aiClient, skill *Skill) (*StyleResult, error) {
	instructions, err := renderTemplate(styleAnalysisPrompt, skill)
	if err != nil {
		return nil, err
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(skill.evidenceText()),
		},
		ResponseFormat: jsonSchemaFormat("style_result", styleResultSchema),
	}

	resp, err := client.Send(params)
	if err != nil {
		return nil, err
	}

	result, err := parse[StyleResult](resp)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type EvidenceSummary struct {
	CompetencyID string `json:"competency_id"`
	Summary      string `json:"summary"`
}

type EvidenceResult struct {
	Competencies []EvidenceSummary `json:"competencies"`
}

var evidenceResultSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"competencies": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"competency_id": map[string]any{"type": "string"},
					"summary":       map[string]any{"type": "string"},
				},
				"required":             []string{"competency_id", "summary"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"competencies"},
	"additionalProperties": false,
}

func AnalyzeEvidence(client aiClient, skill *Skill) (*EvidenceResult, error) {
	instructions, err := renderTemplate(evidenceAnalysisPrompt, skill)
	if err != nil {
		return nil, err
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(skill.evidenceText()),
		},
		ResponseFormat: jsonSchemaFormat("evidence_result", evidenceResultSchema),
	}

	resp, err := client.Send(params)
	if err != nil {
		return nil, err
	}

	result, err := parse[EvidenceResult](resp)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type CoverageItem struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type CoverageResult struct {
	Coverage []CoverageItem `json:"coverage"`
}

var coverageResultSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"coverage": map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{"type": "string"},
					"status": map[string]any{
						"type": "string",
						"enum": []string{"covered", "partial", "none"},
					},
					"reason": map[string]any{"type": "string"},
				},
				"required":             []string{"id", "status", "reason"},
				"additionalProperties": false,
			},
		},
	},
	"required":             []string{"coverage"},
	"additionalProperties": false,
}

func DetectCoverage(client aiClient, skill *Skill) (*CoverageResult, error) {
	instructions, err := renderTemplate(coverageDetectionPrompt, skill)
	if err != nil {
		return nil, err
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(skill.evidenceText()),
		},
		ResponseFormat: jsonSchemaFormat("coverage_result", coverageResultSchema),
	}

	resp, err := client.Send(params)
	if err != nil {
		return nil, err
	}

	result, err := parse[CoverageResult](resp)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

type WriteResult struct {
	Evidence string `json:"evidence"`
}

var writeResultSchema = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"evidence": map[string]any{"type": "string"},
	},
	"required":             []string{"evidence"},
	"additionalProperties": false,
}

func WriteEvidence(client aiClient, skill *Skill, title, review string, style *StyleResult, coverage *CoverageResult, summary *EvidenceResult) (string, error) {
	type writeEvidenceTemplateData struct {
		Category   string
		Competence string
		Title      string
		Style      *StyleResult
		Coverage   *CoverageResult
		Summary    *EvidenceResult
	}

	data := writeEvidenceTemplateData{
		Category:   skill.Category,
		Competence: skill.Competence,
		Title:      title,
		Style:      style,
		Coverage:   coverage,
		Summary:    summary,
	}

	instructions, err := renderTemplate(writeEvidencePrompt, data)
	if err != nil {
		return "", err
	}

	params := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(instructions),
			openai.UserMessage(review),
		},
		ResponseFormat: jsonSchemaFormat("write_result", writeResultSchema),
	}

	resp, err := client.Send(params)
	if err != nil {
		return "", err
	}

	result, err := parse[WriteResult](resp)
	if err != nil {
		return "", err
	}

	return result.Evidence, nil
}
