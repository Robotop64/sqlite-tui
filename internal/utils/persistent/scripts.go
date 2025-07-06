package persistent

import (
	"errors"
	"os"
	"strings"

	utils "github.com/Robotop64/sqlite-tui/internal/utils"
	yaml "gopkg.in/yaml.v3"
)

type MetaData struct {
	Name        string `yaml:"Name"`
	Description string `yaml:"Description"`
}

type Script struct {
	MetaData MetaData
	Script   []byte
}

func LoadScript(path string) (Script, error) {
	path = utils.CleanPath(path)
	linesData, err := os.ReadFile(path)
	if err != nil {
		return Script{}, err
	}
	lines := strings.Split(string(linesData), "\n")
	var start, end int = -1, -1
	for i, c := range lines {
		c = strings.TrimSpace(c)
		if c == "---" {
			if start == -1 {
				start = i
			} else {
				end = i
				break
			}
		}
	}

	if start == -1 || end == -1 || start >= end {
		return Script{}, errors.New("invalid script format: check metadata section")
	}

	metaSection := lines[start+1 : end]
	scriptSection := lines[end+1:]

	var script Script = Script{}

	metaYaml := strings.Join(metaSection, "\n")
	if err := yaml.Unmarshal([]byte(metaYaml), &script.MetaData); err != nil {
		return Script{}, err
	}
	script.Script = []byte(strings.Join(scriptSection, "\n"))

	return script, nil
}
