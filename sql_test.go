package smolid

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestID_Scan(t *testing.T) {
	id := &ID{}

	t.Run("Scan int64", func(t *testing.T) {
		val := int64(123456789)
		err := id.Scan(val)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.n != uint64(val) {
			t.Errorf("expected %v, got %v", val, id.n)
		}
	})

	t.Run("Scan uint64", func(t *testing.T) {
		val := uint64(987654321)
		err := id.Scan(val)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.n != val {
			t.Errorf("expected %v, got %v", val, id.n)
		}
	})

	t.Run("Scan nil", func(t *testing.T) {
		id.n = 123
		err := id.Scan(nil)
		if err != nil {
			t.Fatalf("Scan failed: %v", err)
		}
		if id.n != 0 {
			t.Errorf("expected 0, got %v", id.n)
		}
	})

	t.Run("Scan invalid type", func(t *testing.T) {
		err := id.Scan(true)
		if err == nil {
			t.Error("expected error for bool type, got nil")
		}
	})
}

func TestID_Value(t *testing.T) {
	id := ID{n: 12345}
	val, err := id.Value()
	if err != nil {
		t.Fatalf("Value failed: %v", err)
	}
	if val != int64(12345) {
		t.Errorf("expected int64(12345), got %v (%T)", val, val)
	}

	var _ driver.Valuer = ID{}
	var _ sql.Scanner = &ID{}
}

func TestID_PGX(t *testing.T) {
	id := New()

	t.Run("Int64Value", func(t *testing.T) {
		iv, err := id.Int64Value()
		if err != nil {
			t.Fatalf("Int64Value failed: %v", err)
		}
		if iv.Int64 != int64(id.n) {
			t.Errorf("expected %v, got %v", int64(id.n), iv.Int64)
		}
	})

	t.Run("ScanInt64", func(t *testing.T) {
		other := &ID{}
		iv := pgtype.Int8{Int64: int64(id.n), Valid: true}
		err := other.ScanInt64(iv)
		if err != nil {
			t.Fatalf("ScanInt64 failed: %v", err)
		}
		if other.n != id.n {
			t.Errorf("expected %v, got %v", id.n, other.n)
		}
	})

	t.Run("GormDataType", func(t *testing.T) {
		if id.GormDataType() != "bigint" {
			t.Errorf("expected bigint, got %v", id.GormDataType())
		}
	})
}
