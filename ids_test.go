package smolid

import (
	"math"
	"reflect"
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

func TestID_IsTyped(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want bool
	}{
		{
			name: "Untyped",
			id:   New(),
			want: false,
		},
		{
			name: "Typed",
			id:   Must(NewWithType(1)),
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsTyped(); got != tt.want {
				t.Errorf("IsTyped() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestID_IsOfType(t *testing.T) {
	type args struct {
		typ byte
	}
	tests := []struct {
		name    string
		id      ID
		args    args
		want    bool
		wantErr bool
	}{
		{
			name:    "Invalid Type",
			id:      New(),
			args:    args{typ: v1TypeSize + 1},
			wantErr: true,
		},
		{
			name:    "Match",
			id:      Must(NewWithType(1)),
			args:    args{typ: 1},
			want:    true,
			wantErr: false,
		},
		{
			name: "No Match",
			id:   Must(NewWithType(1)),
			args: args{typ: 2},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.id.IsOfType(tt.args.typ)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsOfType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsOfType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromUint64(t *testing.T) {
	type args struct {
		n uint64
	}

	testID := Must(FromString("apviiguuvmsh2"))

	tests := []struct {
		name    string
		args    args
		want    ID
		wantErr bool
	}{
		{
			name:    "invalid (version missing)",
			args:    args{n: 0},
			wantErr: true,
		},
		{
			name:    "invalid (version too high)",
			args:    args{n: math.MaxUint64},
			wantErr: true,
		},
		{
			name:    "valid (version 1)",
			args:    args{n: testID.n},
			want:    testID,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromUint64(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromUint64() got = %v, want %v", got, tt.want)
			}
		})
	}
}
