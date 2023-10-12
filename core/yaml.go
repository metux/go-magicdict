package core

import (
	"os"

	"gopkg.in/yaml.v3"
)

func YamlParse(text []byte) (*Dict, error) {
	d := EmptyDict()
	err := yaml.Unmarshal(text, &d)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func YamlLoad(fn string) (*Dict, error) {
	text, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	return YamlParse(text)
}
