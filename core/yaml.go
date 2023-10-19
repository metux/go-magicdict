package core

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/metux/go-magicdict/api"
)

// load a Dict from yaml from byte array (using yaml.v3)
func YamlParse(text []byte) (Dict, error) {
	d := EmptyDict()
	err := yaml.Unmarshal(text, &d)
	return d, err
}

// load a Dict from yaml file (using yaml.v3)
func YamlLoad(fn string) (Dict, error) {
	if fn == "" {
		return EmptyDict(), nil
	}
	text, err := os.ReadFile(fn)
	if err != nil {
		return EmptyDict(), err
	}
	return YamlParse(text)
}

// encode an entry (and it sub-entries) into yaml stream (using yaml.v3)
func YamlEncode(root api.Entry) ([]byte, error) {
	return yaml.Marshal(root)
}

// store entry (and it sub-entries) as yaml file (using yaml.v3)
func YamlStore(fn string, root api.Entry, fmode os.FileMode) error {
	if data, err := YamlEncode(root); err != nil {
		return err
	} else {
		return ioutil.WriteFile(fn, data, fmode)
	}
}
