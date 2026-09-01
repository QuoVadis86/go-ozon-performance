package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// Uint64 is a JSON number-or-string tolerant unsigned 64-bit integer.
//
// The Performance API spec declares identifiers, bids and budgets as
// `string` with `uint64` format, but live responses sometimes carry them as
// plain JSON numbers. Uint64 accepts both on decode and always emits the
// canonical string form (matching the official spec) on encode.
type Uint64 uint64

func (u *Uint64) UnmarshalJSON(b []byte) error {
	trimmed := bytes.TrimSpace(b)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		return nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		if s == "" {
			return nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return fmt.Errorf("transport.Uint64: invalid string %q: %w", s, err)
		}
		*u = Uint64(n)
		return nil
	}
	n, err := strconv.ParseUint(string(trimmed), 10, 64)
	if err != nil {
		return fmt.Errorf("transport.Uint64: invalid number %q: %w", string(trimmed), err)
	}
	*u = Uint64(n)
	return nil
}

func (u Uint64) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(strconv.FormatUint(uint64(u), 10))), nil
}

func (u Uint64) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

// U returns the underlying uint64 value.
func (u Uint64) U() uint64 { return uint64(u) }
