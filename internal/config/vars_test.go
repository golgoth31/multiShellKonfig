package config

import (
	"testing"

	"uuid"
)

func TestDefaultConfigIDIsValidUUIDv4(t *testing.T) {
	konfigs := DefaultConfig.Konfigs
	if len(konfigs) == 0 {
		t.Fatal("DefaultConfig has no konfigs")
	}

	id := konfigs[0].ID
	parsed, err := uuid.Parse(id)
	if err != nil {
		t.Fatalf("DefaultConfig ID %q is not a valid UUID: %v", id, err)
	}

	if len(id) != 36 {
		t.Errorf("expected a 36-char RFC 9562 string, got %q (len %d)", id, len(id))
	}
	if parsed[6]>>4 != 4 {
		t.Errorf("expected UUID version 4, got %d (id %q)", parsed[6]>>4, id)
	}
	if parsed[8]>>6 != 0b10 {
		t.Errorf("expected RFC 4122 variant, got %d (id %q)", parsed[8]>>6, id)
	}
}
