package smolid

import (
	"math"
	"strconv"
	"testing"
	"time"
)

var testTimestamp = time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)

func TestRandomSpace(t *testing.T) {
	toFill := uint64(0) | timestampMask | versionMask | v1TypeFlag
	filled := false
	for i := 0; i < 1000; i++ {
		toFill |= New().n
		if toFill == math.MaxUint64 {
			filled = true
			t.Log("Random space filled at " + strconv.Itoa(i))
			break
		}
	}
	if !filled {
		t.Errorf("Random space not filled")
	}
}

func TestNew(t *testing.T) {
	id := New()
	// Check the version
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}
	// Check the type is untyped
	if _, err := id.Type(); err != ErrUntyped {
		t.Errorf("Expected ErrUntyped, got %v", err)
	}
}

func TestIDTimestamp(t *testing.T) {
	now = func() time.Time { return testTimestamp }
	id := New()
	// make sure it has the right timestamp
	if !id.Time().Equal(testTimestamp) {
		t.Errorf("Invalid time. Expected %v, got %v", testTimestamp, id.Time())
	}
	// and the right version
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}
}

func TestIDFromString(t *testing.T) {
	id, err := FromString("ACPJE64AEYEZ6")
	if err != nil {
		t.Fatal(err)
	}
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}
	// same thing for lowercase
	id, err = FromString("acpje64aeyez6")
	if err != nil {
		t.Fatal(err)
	}
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}
}

func TestIDFromStringInvalid(t *testing.T) {
	_, err := FromString("ACPJE64AEYEZ")
	if err == nil {
		t.Fatal("Expected error")
	}
	_, err = FromString("invalid")
	if err == nil {
		t.Fatal("Expected error")
	}
}

func TestNewForType(t *testing.T) {
	now = func() time.Time { return testTimestamp }
	const (
		MyType = iota
		MyOtherType
	)
	id, err := NewWithType(MyType)

	if err != nil {
		t.Fatal(err)
	}

	// make sure it has the right timestamp still
	if !id.Time().Equal(testTimestamp) {
		t.Errorf("Invalid time. Expected %v, got %v", testTimestamp, id.Time())
	}
	// and the right version still
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}

	// and finally make sure we're able to extract the embedded type id
	if typ, _ := id.Type(); typ != MyType {
		t.Errorf("Expected type %v, got %v", MyType, typ)
	}

	// Then the same deal for the other type
	id, err = NewWithType(MyOtherType)
	if err != nil {
		t.Fatal(err)
	}
	if !id.Time().Equal(testTimestamp) {
		t.Errorf("Invalid time. Expected %v, got %v", testTimestamp, id.Time())
	}
	if id.Version() != 1 {
		t.Errorf("Expected version 1, got %v", id.Version())
	}
	if typ, _ := id.Type(); typ != MyOtherType {
		t.Errorf("Expected type %v, got %v", MyOtherType, typ)
	}
}

func TestNewForTypeInvalid(t *testing.T) {
	_, err := NewWithType(v1TypeSize + 1)
	if err == nil {
		t.Fatal("Expected error")
	}
}
