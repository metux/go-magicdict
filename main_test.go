package magicdict

import (
	"testing"

	"github.com/metux/go-magicdict/api"
	"github.com/metux/go-magicdict/core"
	"github.com/metux/go-magicdict/magic"
	"github.com/metux/go-magicdict/tests"
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
	core.YamlStore("test-output1.yaml", root, 0644)
	tests.RunTestOne(t, root)
	core.YamlStore("test-output2.yaml", root, 0644)
}
