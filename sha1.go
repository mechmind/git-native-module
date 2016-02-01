package git

import (
	"fmt"

	"github.com/mechmind/git-go/rawgit"
)

type sha1 rawgit.OID

// Equal returns true if s has the same sha1 as caller.
// Support 40-length-string, []byte, sha1.
func (id sha1) Equal(s2 interface{}) bool {
	switch v := s2.(type) {
	case string:
		if len(v) != 40 {
			return false
		}
		return v == id.String()
	case []byte:
		if len(v) != 20 {
			return false
		}
		for i, v := range v {
			if id[i] != v {
				return false
			}
		}
	case sha1:
		return id == v
	default:
		return false
	}
	return true
}

// String returns string (hex) representation of the Oid.
func (s sha1) String() string {
	return sha2oidp(s).String()
}

// MustID always creates a new sha1 from a [20]byte array with no validation of input.
func MustID(b []byte) sha1 {
	var id sha1
	for i := 0; i < 20; i++ {
		id[i] = b[i]
	}
	return id
}

// NewID creates a new sha1 from a [20]byte array.
func NewID(b []byte) (sha1, error) {
	if len(b) != 20 {
		return sha1{}, fmt.Errorf("Length must be 20: %v", b)
	}
	return MustID(b), nil
}

// MustIDFromString always creates a new sha from a ID with no validation of input.
func MustIDFromString(s string) sha1 {
	oid, _ := rawgit.ParseOID(s)
	return sha1(*oid)
}

// NewIDFromString creates a new sha1 from a ID string of length 40.
func NewIDFromString(s string) (sha1, error) {
	oid, err := rawgit.ParseOID(s)
	if err != nil {
		return sha1{}, err
	}

	return sha1(*oid), nil
}

func sha2oidp(id sha1) *rawgit.OID {
	oid := rawgit.OID(id)
	return &oid
}
