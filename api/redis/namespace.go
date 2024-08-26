package redis

import (
	"strings"

	"github.com/eosswedenorg/thalos/api"
)

const (
	// Default prefix to use when none is set.
	defaultPrefix = "ship"

	// We need to have some chain_id, so if no one is specified.
	// we use a "null" id that is all zeros.
	nullChain = "0000000000000000000000000000000000000000000000000000000000000000"
)

// Namespace type.
//
// Contains a prefix and chain_id to guard keys against collision.
// Prefix should be sufficient to not collide with other application using the same redis database.
// chain_id should be fine to not let multiple reader with different chains to write to the same channels.

type Namespace struct {
	Prefix  string
	ChainID string
}

// Create a new key with this namespace.
func (ns Namespace) NewKey(ch api.Channel) Key {
	return Key{NS: ns, Channel: ch}
}

func (ns Namespace) String() string {
	// No Chain id, set to "nullChain"
	if len(ns.ChainID) < 1 {
		ns.ChainID = nullChain
	}

	// Set default prefix if empty.
	if len(ns.Prefix) < 1 {
		ns.Prefix = defaultPrefix
	}

	// Otherwise. return both.
	return strings.Join([]string{ns.Prefix, ns.ChainID}, "::")
}
