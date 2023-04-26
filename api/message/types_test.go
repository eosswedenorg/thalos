package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypes_ActionTraceGetData(t *testing.T) {
	act := ActionTrace{
		Data: []byte(`{"one": 1234, "two": "string"}`),
	}

	actual, err := act.GetData()
	assert.NoError(t, err)

	exptected := map[string]interface{}{
		"one": float64(1234),
		"two": "string",
	}

	assert.Equal(t, exptected, actual)
}
