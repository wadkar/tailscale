// Copyright (c) 2020 Tailscale Inc & AUTHORS All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package key defines some types related to curve25519 keys.
package key

import (
	"encoding/base64"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

// Private represents a curve25519 private key.
type Private [32]byte

// Private reports whether p is the zero value.
func (p Private) IsZero() bool { return p == Private{} }

// B32 returns k as the *[32]byte type that's used by the
// golang.org/x/crypto packages. This allocates; it might
// not be appropriate for performance-sensitive paths.
func (k Private) B32() *[32]byte { return (*[32]byte)(&k) }

// Public represents a curve25519 public key.
type Public [32]byte

// Public reports whether p is the zero value.
func (p Public) IsZero() bool { return p == Public{} }

// ShortString returns the Tailscale conventional debug representation
// of a public key: the first five base64 digits of the key, in square
// brackets.
func (p Public) ShortString() string {
	return "[" + base64.StdEncoding.EncodeToString(p[:])[:5] + "]"
}

func (p Public) MarshalText() ([]byte, error) {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(p)))
	base64.StdEncoding.Encode(buf, p[:])
	return buf, nil
}

func (p *Public) UnmarshalText(txt []byte) error {
	if *p != (Public{}) {
		return errors.New("refusing to unmarshal into non-zero key.Public")
	}
	n, err := base64.StdEncoding.Decode(p[:], txt)
	if err != nil {
		return err
	}
	if n != 32 {
		return fmt.Errorf("short decode of %d; want 32", n)
	}
	return nil
}

// B32 returns k as the *[32]byte type that's used by the
// golang.org/x/crypto packages. This allocates; it might
// not be appropriate for performance-sensitive paths.
func (k Public) B32() *[32]byte { return (*[32]byte)(&k) }

func (k Private) Public() Public {
	var pub [32]byte
	curve25519.ScalarBaseMult(&pub, (*[32]byte)(&k))
	return Public(pub)
}
