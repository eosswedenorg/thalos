package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlacklist_Empty(t *testing.T) {
	bl := Blacklist{
		table: map[string][]string{},
	}

	require.True(t, bl.Empty())

	bl.Add("contract", "action1")

	require.False(t, bl.Empty())
}

func TestBlacklist_Add(t *testing.T) {
	bl := Blacklist{
		table: map[string][]string{},
	}
	bl.Add("contract", "action1")
	bl.Add("contract", "action2")
	bl.Add("contract2", "action1")

	expected := Blacklist{
		table: map[string][]string{
			"contract":  {"action1", "action2"},
			"contract2": {"action1"},
		},
	}

	require.Equal(t, expected, bl)
}

func TestBlacklist_IsAllowed(t *testing.T) {
	bl := Blacklist{
		table: map[string][]string{
			"mycontract": {"myaction", "noop"},
		},
	}

	require.False(t, bl.IsAllowed("mycontract", "myaction"))
	require.False(t, bl.IsAllowed("mycontract", "noop"))
	require.True(t, bl.IsAllowed("mycontract", "xxx"))
	require.True(t, bl.IsAllowed("xxx", "yyy"))
}

func TestBlacklist_IsAllowedWildcard(t *testing.T) {
	bl := Blacklist{
		table: map[string][]string{
			"mycontract":   {"*"},
			"*":            {"action1", "action2"},
			"evilcontract": {"evilaction"},
		},
	}

	require.False(t, bl.IsAllowed("mycontract", "myaction"))
	require.False(t, bl.IsAllowed("mycontract", "noop"))
	require.False(t, bl.IsAllowed("mycontract", "xxx"))

	// Wildcard contract
	require.False(t, bl.IsAllowed("somecontract", "action1"))
	require.False(t, bl.IsAllowed("someothercontract", "action1"))
	require.False(t, bl.IsAllowed("randomcontract", "action2"))
	require.False(t, bl.IsAllowed("evilcontract", "action2"))
	require.False(t, bl.IsAllowed("evilcontract", "evilaction"))

	require.True(t, bl.IsAllowed("xxx", "yyy"))
	require.True(t, bl.IsAllowed("evilcontract", "alloweaction"))
}

func TestBlacklist_Whitelist(t *testing.T) {
	bl := Blacklist{
		table: map[string][]string{
			"mycontract": {"myaction", "noop"},
			"*":          {"goodaction1", "goodaction2"},
		},
	}

	bl.SetWhitelist(true)

	require.True(t, bl.IsAllowed("mycontract", "myaction"))
	require.True(t, bl.IsAllowed("mycontract", "noop"))

	// Wildcard contract
	require.True(t, bl.IsAllowed("mycontract", "goodaction1"))
	require.True(t, bl.IsAllowed("someothercontract", "goodaction2"))

	require.False(t, bl.IsAllowed("mycontract", "xxx"))
	require.False(t, bl.IsAllowed("xxx", "yyy"))
}
