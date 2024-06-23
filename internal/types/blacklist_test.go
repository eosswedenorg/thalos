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

func TestBlacklist_Lookup(t *testing.T) {
	bl := Blacklist{
		"mycontract": {"myaction", "noop"},
	}

	require.True(t, bl.Lookup("mycontract", "myaction"))
	require.True(t, bl.Lookup("mycontract", "noop"))
	require.False(t, bl.Lookup("mycontract", "xxx"))
	require.False(t, bl.Lookup("xxx", "yyy"))
}

func TestBlacklist_LookupWildcard(t *testing.T) {
	bl := Blacklist{
		"mycontract": {"*"},
	}

	require.True(t, bl.Lookup("mycontract", "myaction"))
	require.True(t, bl.Lookup("mycontract", "noop"))
	require.True(t, bl.Lookup("mycontract", "xxx"))
	require.False(t, bl.Lookup("xxx", "yyy"))
}
