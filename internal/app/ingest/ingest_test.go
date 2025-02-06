package ingest

import (
	"encoding/json"
	"testing"
)

func TestSerialize(t *testing.T) {
	data := map[string][]string{
		"key1": {"value1", "value2"},
		"key2": {"value3", "value4"},
	}
	expected, _ := json.Marshal(data)

	result, err := serialize(data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(result) != string(expected) {
		t.Errorf("Expected %s, got %s", string(expected), string(result))
	}
}

func TestDeserialize(t *testing.T) {
	data := map[string][]string{
		"key1": {"value1", "value2"},
		"key2": {"value3", "value4"},
	}
	serialized, _ := json.Marshal(data)
	result, err := deserialize(serialized)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(result) != len(data) {
		t.Errorf("Expected %v, got %v", len(data), len(result))
	}
	for k, v := range data {
		if len(result[k]) != len(v) {
			t.Errorf("Expected %v, got %v", len(v), len(result[k]))
		}
		for i, j := range v {
			if j != result[k][i] {
				t.Errorf("Expected %v, got %v", j, result[k][i])
			}
		}
	}
}
