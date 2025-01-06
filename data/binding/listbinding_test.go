package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type simpleList struct {
	listBase
}

func TestListBase_AddListener(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, data.listeners.Len())

	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data.AddListener(fn)
	assert.Equal(t, 1, data.listeners.Len())

	data.trigger()
	assert.True(t, called)
}

func TestListBase_GetItem(t *testing.T) {
	data := &simpleList{}
	f := 0.5
	data.appendItem(BindFloat(&f))
	assert.Len(t, data.items, 1)

	item, err := data.GetItem(0)
	require.NoError(t, err)
	val, err := item.(Float).Get()
	require.NoError(t, err)
	assert.Equal(t, f, val)

	_, err = data.GetItem(5)
	assert.Error(t, err)
}

func TestListBase_Length(t *testing.T) {
	data := &simpleList{}
	assert.Equal(t, 0, data.Length())

	data.appendItem(NewFloat())
	assert.Equal(t, 1, data.Length())
}

func TestListBase_RemoveListener(t *testing.T) {
	called := false
	fn := NewDataListener(func() {
		called = true
	})
	data := &simpleList{}
	data.listeners.Store(fn, true)

	assert.Equal(t, 1, data.listeners.Len())
	data.RemoveListener(fn)
	assert.Equal(t, 0, data.listeners.Len())

	data.trigger()
	assert.False(t, called)
}
