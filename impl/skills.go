package impl

import (
	"os"

	"github.com/gocarina/gocsv"
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
