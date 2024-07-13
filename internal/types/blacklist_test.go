package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlacklist_Add(t *testing.T) {
	bl := Blacklist{}
	bl.Add("contract", "action1")
	bl.Add("contract", "action2")
	bl.Add("contract2", "action1")

	expected := Blacklist{
		"contract":  {"action1", "action2"},
		"contract2": {"action1"},
	}

	require.Equal(t, expected, bl)
}

func TestBlacklist_IsAllowed(t *testing.T) {
	bl := Blacklist{
		"mycontract": {"myaction", "noop"},
	}

	require.False(t, bl.IsAllowed("mycontract", "myaction"))
	require.False(t, bl.IsAllowed("mycontract", "noop"))
	require.True(t, bl.IsAllowed("mycontract", "xxx"))
	require.True(t, bl.IsAllowed("xxx", "yyy"))
}

func TestBlacklist_IsAllowedWildcard(t *testing.T) {
	bl := Blacklist{
		"mycontract": {"*"},
	}

	require.False(t, bl.IsAllowed("mycontract", "myaction"))
	require.False(t, bl.IsAllowed("mycontract", "noop"))
	require.False(t, bl.IsAllowed("mycontract", "xxx"))
	require.True(t, bl.IsAllowed("xxx", "yyy"))
}
