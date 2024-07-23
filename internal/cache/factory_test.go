package cache_test

import (
	"testing"

	"github.com/eosswedenorg/thalos/internal/cache"
	"github.com/stretchr/testify/require"
)

func TestFactory_Make(t *testing.T) {
	store, err := cache.Make("memory", map[string]any{})
	require.NoError(t, err)
	require.Equal(t, cache.NewMemoryStore(), store)
}

func TestFactory_MakeInvalidDriver(t *testing.T) {
	store, err := cache.Make("87923yus", map[string]any{})
	require.Error(t, err)
	require.Nil(t, store)
}
