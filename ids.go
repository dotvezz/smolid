package smolid

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"strings"
	"time"
)

// Taking a reference to time.Now to make testing easier
var now = time.Now

const (
	epoch                uint64 = 1735707600000 // 2025-01-01 00:00:00
	timestampSize               = 0b11111111111111111111111111111111111111111
	timestampShiftOffset        = 23 // The timestamp is shifted three bytes to the left
	timestampMask               = timestampSize << timestampShiftOffset
	versionShiftOffset          = 21
	versionMask                 = 0b11 << versionShiftOffset

	v1TypeShiftOffset = 9
	v1TypeFlag        = 0b1 << 20
	v1TypeSize        = 0b1111111
	v1TypeMask        = v1TypeSize << v1TypeShiftOffset
	v1RandomSpace     = 0xfffff // In V1, the least significant two and a half bytes (20 bits) can be random
	v1Version         = 0b1 << 21
)

/*
ID is a 64-bit (8-byte) value intended to be
  - URL-Friendly; short and unobtrusive in its default unpadded base32 string encoding
  - temporally sortable with strong index locality
  - fast-enough and unique-enough for most use cases

Field Definitions

	0                   1                   2                   3
	0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|                          time_high                            |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
	|    time_low     |ver|t| rand  | type or rand|       rand      |
	+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

Field Descriptions:

  - Timestamp (41 bits): The most significant 41 bits represent a millisecond-precision timestamp
    The allowed timestamp range is 2025-01-01 00:00:00 - 2094-09-07 15:47:35
  - Version (2 bits): Bits 41-42 are reserved for versioning.
    v1 is `01`
  - Type Flag (1 bit): Bit 43 serves as a boolean flag. If set, the "Type/Rand" field is an embedded type identifier.
  - Random (4 bits): The remaining 4 bits of the 6th byte are populated with pseudo-random data.
  - Type/Random (7 bits): If the Type Flag is set, this field contains the Type Identifier. Otherwise, it
    is populated with pseudo-random data.
  - Random (9 bits): The remaining byte is dedicated to pseudo-random data to _reasonably_ ensure uniqueness.

String Format

	The string format is base32 with no padding. Canonically the string is lowercased. This decision is purely for
	aesthetics, but the parser is case-insensitive and will accept uppercase base32 strings.
*/
type ID struct {
	n uint64
}

// New returns a new smolid.ID v1 with all defaults.
func New() ID {
	var id = (uint64(now().UTC().UnixMilli()) - epoch) << timestampShiftOffset // set the timestamp
	id |= v1Version                                                            // set the version bit
	id |= rand.Uint64N(v1RandomSpace)                                          // radom-fill the remaining space
	return ID{id}
}

func Nil() ID { return ID{0} }

// NewWithType returns a new smolid.ID v1 with the given type identifier embedded into the ID.
func NewWithType(typ byte) (ID, error) {
	if typ > v1TypeSize {
		return Nil(), ErrInvalidType
	}
	id := New()                              // get a new v1 ID
	id.n &^= v1TypeMask                      // clear the random data in the type space
	id.n |= v1TypeFlag                       // set the type flag
	id.n |= uint64(typ) << v1TypeShiftOffset // set the type
	return id, nil
}

// FromString parses a smolid.ID from a string. While the canonical representation is all-lowercase, the parser is
// case-insensitive and will accept uppercase or mixed case without problems.
func FromString(s string) (_ ID, err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()

	s = strings.ToUpper(s)
	var bs []byte
	bs, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(s)
	if err != nil {
		return Nil(), err
	}

	return ID{binary.BigEndian.Uint64(bs)}, nil
}

// Must is a convenience function that panics if the given error is not nil. Useful for testing or scenarios where you
// know fully that an ID is valid, but don't want to deal with ignoring the error.
//
// Stolen from gofrs/uuid
func Must(id ID, err error) ID {
	if err != nil {
		panic("couldn't parse id: " + err.Error())
	}
	return id
}

// String returns the canonical string representation of the ID.
func (id ID) String() string {
	return strings.ToLower(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(id.Bytes()))
}

// Version returns the version of the ID
func (id ID) Version() int { return int(id.n&versionMask) >> versionShiftOffset }

// Bytes returns the raw 64-bit integer representation of the ID as a []byte.
func (id ID) Bytes() []byte {
	var bs = make([]byte, 8)
	binary.BigEndian.PutUint64(bs, id.n)
	return bs
}

// Time returns the time.Time embedded in the ID, with millisecond precision.
func (id ID) Time() time.Time {
	var ms = int64(id.n>>timestampShiftOffset + epoch)           // extract the timestamp
	return time.Unix(ms/1000, (ms%1000)*int64(time.Millisecond)) // fill into a time.Time
}

// Type returns the type identifier embedded in the ID, if any. It will return an error if the ID was not created by
// NewWithType.
func (id ID) Type() (byte, error) {
	if !id.IsTyped() {
		return 0, ErrUntyped
	}
	typ := id.n & v1TypeMask
	return byte(typ >> v1TypeShiftOffset), nil
}

// IsTyped returns true if the ID was created by NewWithType.
func (id ID) IsTyped() bool { return id.n&v1TypeFlag != 0 }

// IsOfType returns true if the ID is typed and matches the given type identifier.It will return an error if the ID
// was not created by NewWithType.
func (id ID) IsOfType(typ byte) (bool, error) {
	if !id.IsTyped() {
		return false, ErrUntyped
	}
	if typ > v1TypeSize {
		return false, ErrInvalidType
	}

	// If we've reached this point, we know the ID is typed.
	typ2, _ := id.Type()

	return typ == typ2, nil
}

// Uint64 returns the raw 64-bit integer representation of the ID. Use this instead of Int64Value for most cases. The
// Int64Value method is provided for compatibility with the pgtype.Int8Valuer interface.
func (id ID) Uint64() uint64 { return id.n }
