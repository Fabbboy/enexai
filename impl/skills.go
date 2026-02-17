package impl

import (
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const (
	describeSkillSystem string = `
SYSTEM PROMPT: PREFLIGHT COMPETENCY ANALYZER

You analyze a competency document and optional evidence text.
You do not write final competency evidence.
You only extract structure and guidance.

INPUT
One text block containing:

* competency definitions with IDs like A1.1 and optional K-levels like (K3)
* optional user evidence text
* evidence may reference competencies like (A1.1, A1.3)

TASKS

1. Extract competencies
   For each competency output:
   ID | K-level | short description
   If K-level is missing write K?.

2. Detect coverage from the evidence text
   covered = concrete action described
   partial = mentioned without concrete action
   none = not mentioned
   Always output one line per competency with a short factual reason.

3. Extract writing style from the evidence text
   Always output values. If unknown write unknown.
   Language
   Perspective (first person singular, plural, impersonal)
   Tense (past, present, mixed)
   Tone (formal, semi-formal, narrative, bullet-like)
   Sentence length (short, medium, long, mixed)
   References (yes/no)

4. List gaps
   For competencies with partial or none output:
   ID | very short description of missing evidence type

OUTPUT FORMAT (MARKDOWN)

## COMPETENCIES

A1.1 | K3 | description

## COVERAGE

A1.1 | covered | reason

## STYLE

Language: value
Perspective: value
Tense: value
Tone: value
Sentence length: value
References: yes/no

## GAPS

A1.4 | missing evidence type

CONSTRAINTS
Do not write final competency text.
Do not invent actions.
Be extractive and concise.
One line per competency.

END
`
)

type Skill struct {
	Category   string `csv:"Kompetenzkategorie"`
	Competence string `csv:"Kompetenz"`
	Note       string `csv:"Bemerkung"`
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

type DescribeSkillData struct {
	ResponseDataBase
	skill *Skill
}

func DescribeSkill(data DescribeSkillData) (string, error) {
	var sb strings.Builder
	sb.WriteString("Category: " + data.skill.Category + "\n")
	sb.WriteString("Competence: " + data.skill.Competence + "\n")
	sb.WriteString("Note: " + data.skill.Note + "\n")

	prompt := sb.String()
	input := responses.ResponseNewParamsInputUnion{
		OfString: openai.String(prompt),
	}
	respData := ResponseData{
		ResponseDataBase: data.ResponseDataBase,
		system:           openai.String(describeSkillSystem),
		input:            input,
	}

	resp, err := Response(respData)
	if err != nil {
		return "", err
	}

	return resp, nil
}
