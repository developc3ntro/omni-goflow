package types_test

import (
	"encoding/json"
	"testing"

	"github.com/nyaruka/goflow/excellent/types"

	"github.com/stretchr/testify/assert"
)

func TestXText(t *testing.T) {
	// test equality
	assert.True(t, types.NewXText("abc").Equals(types.NewXText("abc")))
	assert.False(t, types.NewXText("abc").Equals(types.NewXText("def")))

	// test comparison
	assert.Equal(t, 0, types.NewXText("abc").Compare(types.NewXText("abc")))
	assert.Equal(t, 1, types.NewXText("def").Compare(types.NewXText("abc")))
	assert.Equal(t, -1, types.NewXText("abc").Compare(types.NewXText("def")))

	// test length
	assert.Equal(t, 0, types.NewXText("").Length())
	assert.Equal(t, 3, types.NewXText("abc").Length())
	assert.Equal(t, 2, types.NewXText("世界").Length())
	assert.Equal(t, 1, types.NewXText("😁").Length())

	assert.Equal(t, "abc", types.NewXText("abc").String())

	// unmarshal
	var val types.XText
	err := json.Unmarshal([]byte(`"hello"`), &val)
	assert.NoError(t, err)
	assert.Equal(t, types.NewXText("hello"), val)
}
