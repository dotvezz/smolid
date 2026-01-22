# smolid

`smolid` is a 64-bit (8-byte) ID scheme for Go that is URL-friendly, temporally sortable, and optimized for database locality. It is designed for use cases where 128-bit UUIDs are unnecessarily large, but a standard auto-incrementing integer is insufficient.

- **URL-Friendly**: Encoded as short, unpadded base32 strings (e.g., `acpje64aeyez6`).
- **Temporally Sortable**: Most significant bits contain a millisecond-precision timestamp.
- **Compact**: Fits into a standard Go `uint64` and a PostgreSQL `bigint`.
- **Type-Aware**: Supports embedding an optional 7-bit type identifier directly into the ID.

## Related Links

- [mirorac/smolid-js](https://github.com/mirorac/smolid-js) is a reimplementation of `smolid` in Javascript and Typescript by [@mirorac](https://github.com/mirorac)
  - NPM Link: https://www.npmjs.com/package/smolid

## ID Structure

A `smolid` consists of 64 bits partitioned as follows:

```text
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                          time_high                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|    time_low     |ver|t| rand  | type or rand|       rand      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

*   **Timestamp (41 bits)**: Millisecond-precision timestamp with a custom epoch (2025-01-01). Valid until 2094.
*   **Version (2 bits)**: Reserved for versioning (v1 is `01`).
*   **Type Flag (1 bit)**: Boolean flag indicating if the Type field is used.
*   **Random/Type (20 bits)**:
    *   If Type Flag is **unset**: 20 bits of pseudo-random data.
    *   If Type Flag is **set**: 4 bits of random data, a 7-bit Type ID, and 9 bits of random data.


## Usage

The [godoc](https://pkg.go.dev/github.com/dotvezz/smolid?utm_source=godoc) is available as a reference.

### Basic Example

```go
import "github.com/dotvezz/smolid"

// Generate a new ID
id := smolid.New()
fmt.Println(id.String()) // e.g., "acpje64aeyez6"

// Parse from string
parsed, _ := smolid.FromString("acpje64aeyez6")
```

### Using Embedded Types

Embedded types allow you to identify the resource type (e.g., User, Post, Comment) directly from the ID itself.

```go
const (
    TypeUser byte = iota + 1
    TypePost
)

// Create an ID with a type
id, _ := smolid.NewWithType(TypeUser)

// Check type later
if t, err := id.Type(); err == nil {
    switch t {
    case TypeUser:
        fmt.Println("This is a User ID")
    }
}
```

## Integration

`smolid` implements several standard Go interfaces for seamless integration:

- **JSON**: `json.Marshaler` and `json.Unmarshaler`.
- **Text**: `encoding.TextMarshaler` and `encoding.TextUnmarshaler`.
- **SQL**: `database/sql.Scanner` and `database/sql/driver.Valuer`.
- **Postgres**: Native support for `pgx` (via `pgtype.Int8Scanner/Valuer`).
- **GORM**: Automatically identifies as `bigint`.

## Considerations

### Uniqueness and Collisions
`smolid` provides 13 to 20 bits of entropy per millisecond. This is "unique-enough" for many applications but is not a replacement for UUIDs in high-concurrency environments with massive write volumes (e.g., >1000 IDs per millisecond).

### Database Compatibility
While `smolid` fits in a `uint64`, PostgreSQL `bigint` columns are signed. The timestamp will continue to work correctly and remain sortable until the year 2059, at which point the most significant bit flips and the values become negative from a signed perspective.

## License
MIT
