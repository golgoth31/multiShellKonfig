package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golgoth31/multiShellKonfig/internal/config"
	"gopkg.in/yaml.v3"
	"uuid"
)

const fixtureKubeconfig = `apiVersion: v1
kind: Config
current-context: test
`

func TestAddCmdGeneratesUUIDv4(t *testing.T) {
	homeDir := t.TempDir()
	cfgPath := filepath.Join(homeDir, "config.yaml")
	kubeconfigPath := filepath.Join(homeDir, "kubeconfig.yaml")

	if err := os.WriteFile(kubeconfigPath, []byte(fixtureKubeconfig), 0600); err != nil {
		t.Fatalf("write fixture kubeconfig: %v", err)
	}

	existingID := "existing-id"
	homedir = homeDir
	cfgFile = cfgPath
	cfgData = config.Konfigs{
		Konfigs: []config.Konfig{{Path: "~/.kube/config", ID: existingID}},
	}

	addCmd.Run(addCmd, []string{kubeconfigPath})

	data, err := os.ReadFile(cfgFile)
	if err != nil {
		t.Fatalf("read written config: %v", err)
	}

	var got config.Konfigs
	if err := yaml.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal written config: %v", err)
	}

	if len(got.Konfigs) != 2 {
		t.Fatalf("expected 2 konfigs, got %d", len(got.Konfigs))
	}

	added := got.Konfigs[1]
	if added.Path != "~"+"/"+filepath.Base(kubeconfigPath) {
		t.Errorf("expected path to be normalized to ~/%s, got %q", filepath.Base(kubeconfigPath), added.Path)
	}
	if added.ID == existingID {
		t.Errorf("new ID must differ from pre-existing ID %q", existingID)
	}

	parsed, err := uuid.Parse(added.ID)
	if err != nil {
		t.Fatalf("added ID %q is not a valid UUID: %v", added.ID, err)
	}
	if len(added.ID) != 36 {
		t.Errorf("expected a 36-char RFC 9562 string, got %q (len %d)", added.ID, len(added.ID))
	}
	if parsed[6]>>4 != 4 {
		t.Errorf("expected UUID version 4, got %d (id %q)", parsed[6]>>4, added.ID)
	}
	if parsed[8]>>6 != 0b10 {
		t.Errorf("expected RFC 4122 variant, got %d (id %q)", parsed[8]>>6, added.ID)
	}
}

func TestAddCmdDoesNotDuplicatePath(t *testing.T) {
	homeDir := t.TempDir()
	cfgPath := filepath.Join(homeDir, "config.yaml")
	kubeconfigPath := filepath.Join(homeDir, "kubeconfig.yaml")

	if err := os.WriteFile(kubeconfigPath, []byte(fixtureKubeconfig), 0600); err != nil {
		t.Fatalf("write fixture kubeconfig: %v", err)
	}

	homedir = homeDir
	cfgFile = cfgPath
	cfgData = config.Konfigs{}

	addCmd.Run(addCmd, []string{kubeconfigPath})

	var first config.Konfigs
	if err := yaml.Unmarshal(mustReadFile(t, cfgFile), &first); err != nil {
		t.Fatalf("unmarshal written config: %v", err)
	}

	// Reload from disk so the second run sees the entry already present.
	cfgData = first
	addCmd.Run(addCmd, []string{kubeconfigPath})

	var second config.Konfigs
	if err := yaml.Unmarshal(mustReadFile(t, cfgFile), &second); err != nil {
		t.Fatalf("unmarshal written config: %v", err)
	}

	if len(second.Konfigs) != 1 {
		t.Errorf("expected still 1 konfig after duplicate add, got %d", len(second.Konfigs))
	}
	if second.Konfigs[0].ID != first.Konfigs[0].ID {
		t.Errorf("ID changed on duplicate add: %q -> %q", first.Konfigs[0].ID, second.Konfigs[0].ID)
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
