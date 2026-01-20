package smolid

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

var (
	_ driver.Valuer       = ID{}
	_ sql.Scanner         = &ID{}
	_ pgtype.Int64Valuer  = ID{}
	_ pgtype.Int64Scanner = &ID{}
)

// Scan implements the sql.Scanner interface.
func (id *ID) Scan(value any) error {
	if value == nil {
		id.n = 0
		return nil
	}

	switch v := value.(type) {
	case int64:
		id.n = uint64(v)
	case uint64:
		id.n = v
	default:
		return fmt.Errorf("can't scan %T into smolid.ID", value)
	}

	return nil
}

// Value implements the driver.Valuer interface.
func (id ID) Value() (driver.Value, error) {
	return int64(id.n), nil
}

// Int64Value implements the pgtype.Int8Valuer interface (PostgreSQL BIGINT).
func (id ID) Int64Value() (pgtype.Int8, error) {
	return pgtype.Int8{Int64: int64(id.n), Valid: true}, nil
}

// ScanInt64 implements the pgtype.Int8Scanner interface.
func (id *ID) ScanInt64(v pgtype.Int8) error {
	if !v.Valid {
		id.n = 0
		return nil
	}
	id.n = uint64(v.Int64)
	return nil
}

// GormDataType returns the data type for GORM.
func (ID) GormDataType() string {
	return "bigint"
}
