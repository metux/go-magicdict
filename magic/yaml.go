package magic

import (
	"github.com/metux/go-magicdict/core"
)

// Implements yaml.Marshaler interface
func (this MagicDict) MarshalYAML() (interface{}, error) {
	return this.Data, nil
}

// Implements yaml.Marshaler interface
func (this magicScalar) MarshalYAML() (interface{}, error) {
	return this.Data, nil
}

// Load dict and defaults from yaml files
// if a file name is empty, no file is loaded, but empty Dict used instead
func YamlLoad(src string, dflt string) (MagicDict, error) {
	srcDict, err1 := core.YamlLoad(src)
	defDict, err2 := core.YamlLoad(dflt)
	if err1 != nil {
		return NewMagicFromDict(srcDict, defDict), err1
	}
	return NewMagicFromDict(srcDict, defDict), err2
}
