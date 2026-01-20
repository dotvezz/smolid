package smolid

import (
	"encoding/json"
	"testing"

	goccyyaml "github.com/goccy/go-yaml"
	"gopkg.in/yaml.v3"
)

func TestIDMarshalJSON(t *testing.T) {
	id, err := FromString("ACPJE64AEYEZ6")
	if err != nil {
		t.Fatal(err)
	}

	data, err := json.Marshal(id)
	if err != nil {
		t.Fatal(err)
	}

	var expected = `"acpje64aeyez6"`
	if string(data) != expected {
		t.Errorf("Expected %s, got %s", expected, string(data))
	}
}

func TestIDUnmarshalJSON(t *testing.T) {
	var id ID
	err := json.Unmarshal([]byte(`"ACPJE64AEYEZ6"`), &id)
	if err != nil {
		t.Fatal(err)
	}

	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}

	expected := "acpje64aeyez6"
	if id.String() != expected {
		t.Errorf("Expected %s, got %s", expected, id.String())
	}

	// Test lowercase
	var id2 ID
	err = json.Unmarshal([]byte(`"acpje64aeyez6"`), &id2)
	if err != nil {
		t.Fatal(err)
	}

	if id2.String() != expected {
		t.Errorf("Expected %s, got %s", expected, id2.String())
	}
}

func TestIDUnmarshalJSONInvalid(t *testing.T) {
	var id ID
	err := json.Unmarshal([]byte(`"invalid"`), &id)
	if err == nil {
		t.Fatal("Expected error")
	}

	err = json.Unmarshal([]byte(`123`), &id)
	if err == nil {
		t.Fatal("Expected error")
	}

	err = json.Unmarshal([]byte(`null`), &id)
	if err == nil {
		t.Log("Note: unmarshaling null into ID resulted in Nil ID (0)")
	}
}

func TestIDYAML(t *testing.T) {
	id, _ := FromString("ACPJE64AEYEZ6")

	t.Run("gopkg.in/yaml.v3", func(t *testing.T) {
		data, err := yaml.Marshal(id)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "acpje64aeyez6\n" {
			t.Errorf("Expected acpje64aeyez6, got %s", string(data))
		}

		var id2 ID
		err = yaml.Unmarshal(data, &id2)
		if err != nil {
			t.Fatal(err)
		}
		if id2 != id {
			t.Errorf("Expected %v, got %v", id, id2)
		}
	})

	t.Run("github.com/goccy/go-yaml", func(t *testing.T) {
		data, err := goccyyaml.Marshal(id)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "acpje64aeyez6\n" {
			t.Errorf("Expected acpje64aeyez6, got %s", string(data))
		}

		var id2 ID
		err = goccyyaml.Unmarshal(data, &id2)
		if err != nil {
			t.Fatal(err)
		}
		if id2 != id {
			t.Errorf("Expected %v, got %v", id, id2)
		}
	})
}
