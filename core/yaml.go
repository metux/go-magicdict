package core

import (
    "os"
    "gopkg.in/yaml.v3"
    "github.com/metux/go-magicdict/api"
)

func YamlParse(text [] byte) (*Dict, error) {
    m := new(api.AnyMap)
    err := yaml.Unmarshal(text, &m)
    if err != nil {
        return nil, err
    }
    d := NewDict(m)
    return &d, nil
}

func YamlLoad(fn string) (*Dict, error) {
    text, err := os.ReadFile(fn)
    if err != nil {
        return nil, err
    }
    return YamlParse(text)
}
