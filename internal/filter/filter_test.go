package filter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList_Empty(t *testing.T) {
	list := List{
		table: map[string][]string{},
	}

	require.True(t, list.Empty())

	list.Add("contract", "action1")

	require.False(t, list.Empty())
}

func TestList_Add(t *testing.T) {
	list := List{
		table: map[string][]string{},
	}
	list.Add("contract", "action1")
	list.Add("contract", "action2")
	list.Add("contract2", "action1")

	expected := List{
		table: map[string][]string{
			"contract":  {"action1", "action2"},
			"contract2": {"action1"},
		},
	}

	require.Equal(t, expected, list)
}

func TestList_IsIncluded_ExcludeMode(t *testing.T) {
	list := List{
		table: map[string][]string{
			"mycontract": {"myaction", "noop"},
		},
		mode: Exclude,
	}

	require.False(t, list.IsIncluded("mycontract", "myaction"))
	require.False(t, list.IsIncluded("mycontract", "noop"))
	require.True(t, list.IsIncluded("mycontract", "xxx"))
	require.True(t, list.IsIncluded("xxx", "yyy"))
}

func TestList_IsIncludedWildcard_ExcludeMode(t *testing.T) {
	list := List{
		table: map[string][]string{
			"mycontract":   {"*"},
			"*":            {"action1", "action2"},
			"evilcontract": {"evilaction"},
		},
		mode: Exclude,
	}

	require.False(t, list.IsIncluded("mycontract", "myaction"))
	require.False(t, list.IsIncluded("mycontract", "noop"))
	require.False(t, list.IsIncluded("mycontract", "xxx"))

	// Wildcard contract
	require.False(t, list.IsIncluded("somecontract", "action1"))
	require.False(t, list.IsIncluded("someothercontract", "action1"))
	require.False(t, list.IsIncluded("randomcontract", "action2"))
	require.False(t, list.IsIncluded("evilcontract", "action2"))
	require.False(t, list.IsIncluded("evilcontract", "evilaction"))

	require.True(t, list.IsIncluded("xxx", "yyy"))
	require.True(t, list.IsIncluded("evilcontract", "alloweaction"))
}

func TestList_IncludeMode(t *testing.T) {
	list := List{
		table: map[string][]string{
			"mycontract": {"myaction", "noop"},
			"*":          {"goodaction1", "goodaction2"},
		},
		mode: Include,
	}

	require.True(t, list.IsIncluded("mycontract", "myaction"))
	require.True(t, list.IsIncluded("mycontract", "noop"))

	// Wildcard contract
	require.True(t, list.IsIncluded("mycontract", "goodaction1"))
	require.True(t, list.IsIncluded("someothercontract", "goodaction2"))

	require.False(t, list.IsIncluded("mycontract", "xxx"))
	require.False(t, list.IsIncluded("xxx", "yyy"))
}

func TestList_IsExcluded(t *testing.T) {
	list := List{
		table: map[string][]string{
			"mycontract": {"myaction", "noop"},
		},
		mode: Exclude,
	}

	require.True(t, list.IsExcluded("mycontract", "myaction"))
	require.False(t, list.IsExcluded("mycontract", "randomaction"))

	list.SetMode(Include)

	require.False(t, list.IsExcluded("mycontract", "myaction"))
	require.True(t, list.IsExcluded("mycontract", "randomaction"))
}
