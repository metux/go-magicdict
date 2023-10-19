package magicdict

import (
	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
	"github.com/metux/go-magicdict/magic"
	"github.com/metux/go-magicdict/tests"
	"testing"
)

func loadOne(t *testing.T) api.Entry {
	root, err := core.YamlLoad("tests/one.yaml")
	if err != nil {
		t.Fatalf("failed loading yaml: %s", err)
	}

	dflt, err := core.YamlLoad("tests/defaults.yml")
	if err != nil {
		t.Fatalf("failed loading yaml: %s", err)
	}

	return magic.NewMagicFromDict(root, dflt)
}

func loadDefaults(t *testing.T) api.Entry {
	root, err := core.YamlLoad("tests/one.yaml")
	if err != nil {
		t.Fatalf("failed loading yaml: %s", err)
	}
	return root
}

func TestLoadOnly(t *testing.T) {
	loadOne(t)
}

func TestYamlOne(t *testing.T) {
	root := loadOne(t)

	tests.RunTestOne(t, root)
}
