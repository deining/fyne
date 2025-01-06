package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindStringTree(t *testing.T) {
	ids := map[string][]string{DataTreeRootID: {"1", "5", "2"}}
	l := map[string]string{"1": "one", "5": "five", "2": "two and a half"}
	f := BindStringTree(&ids, &l)

	assert.Len(t, f.ChildIDs(DataTreeRootID), 3)
	v, err := f.GetValue("5")
	assert.NoError(t, err)
	assert.Equal(t, "five", v)

	assert.NotNil(t, f.(*boundStringTree).val)
	assert.Len(t, *(f.(*boundStringTree).val), 3)

	_, err = f.GetValue("nan")
	assert.Error(t, err)
}

func TestExternalFloatTree_Reload(t *testing.T) {
	i := map[string][]string{"": {"1", "2"}, "1": {"3"}}
	m := map[string]float64{"1": 1.0, "2": 5.0, "3": 2.3}
	f := BindFloatTree(&i, &m)

	assert.Len(t, f.ChildIDs(""), 2)
	v, err := f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 5.0, v)

	calledTree, calledChild := false, false
	f.AddListener(NewDataListener(func() {
		calledTree = true
	}))
	assert.True(t, calledTree)

	child, err := f.GetItem("2")
	assert.NoError(t, err)
	child.AddListener(NewDataListener(func() {
		calledChild = true
	}))
	assert.True(t, calledChild)

	assert.NotNil(t, f.(*boundFloatTree).val)
	assert.Len(t, *(f.(*boundFloatTree).val), 3)

	_, err = f.GetValue("-1")
	assert.Error(t, err)

	calledTree, calledChild = false, false
	m["2"] = 4.8
	f.Reload()
	v, err = f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 4.8, v)
	assert.False(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2}
	f.Reload()
	v, err = f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledTree)
	assert.True(t, calledChild)

	calledTree, calledChild = false, false
	m = map[string]float64{"1": 1.0, "2": 4.2, "3": 5.3}
	f.Reload()
	v, err = f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 4.2, v)
	assert.True(t, calledTree)
	assert.False(t, calledChild)
}

func TestNewStringTree(t *testing.T) {
	f := NewStringTree()
	assert.Len(t, f.ChildIDs(DataTreeRootID), 0)

	_, err := f.GetValue("NaN")
	assert.Error(t, err)
}

func TestStringTree_Append(t *testing.T) {
	f := NewStringTree()
	assert.Len(t, f.ChildIDs(DataTreeRootID), 0)

	f.Append(DataTreeRootID, "5", "five")
	assert.Len(t, f.ChildIDs(DataTreeRootID), 1)
}

func TestStringTree_GetValue(t *testing.T) {
	f := NewStringTree()

	err := f.Append(DataTreeRootID, "1", "1.3")
	assert.NoError(t, err)
	v, err := f.GetValue("1")
	assert.NoError(t, err)
	assert.Equal(t, "1.3", v)

	err = f.Append(DataTreeRootID, "fraction", "0.2")
	assert.NoError(t, err)
	v, err = f.GetValue("fraction")
	assert.NoError(t, err)
	assert.Equal(t, "0.2", v)

	err = f.SetValue("1", "0.5")
	assert.NoError(t, err)
	v, err = f.GetValue("1")
	assert.NoError(t, err)
	assert.Equal(t, "0.5", v)
}

func TestStringTree_Remove(t *testing.T) {
	f := NewStringTree()
	f.Append(DataTreeRootID, "5", "five")
	f.Append(DataTreeRootID, "3", "three")
	f.Append("5", "53", "fifty three")
	assert.Len(t, f.ChildIDs(DataTreeRootID), 2)
	assert.Len(t, f.ChildIDs("5"), 1)

	f.Remove("5")
	assert.Len(t, f.ChildIDs(DataTreeRootID), 1)
	assert.Len(t, f.ChildIDs("5"), 0)
}

func TestFloatTree_Set(t *testing.T) {
	ids := map[string][]string{"": {"1", "2"}, "1": {"3"}}
	m := map[string]float64{"1": 1.0, "2": 5.0, "3": 2.3}
	f := BindFloatTree(&ids, &m)
	i, err := f.GetItem("2")
	assert.NoError(t, err)
	data := i.(Float)

	assert.Len(t, f.ChildIDs(""), 2)
	v, err := f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 5.0, v)
	v, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5.0, v)

	ids = map[string][]string{"": {"1", "2"}, "1": {"3", "4"}}
	m = map[string]float64{"1": 1.2, "2": 5.2, "3": 2.2, "4": 4.2}
	err = f.Set(ids, m)
	assert.NoError(t, err)

	assert.Len(t, f.ChildIDs("1"), 2)
	v, err = f.GetValue("2")
	assert.NoError(t, err)
	assert.Equal(t, 5.2, v)
	v, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5.2, v)

	ids = map[string][]string{"": {"1", "2"}}
	m = map[string]float64{"1": 1.3, "2": 5.3}
	err = f.Set(ids, m)
	assert.NoError(t, err)

	assert.Len(t, f.ChildIDs(""), 2)
	v, err = f.GetValue("1")
	assert.NoError(t, err)
	assert.Equal(t, 1.3, v)
	v, err = data.Get()
	assert.NoError(t, err)
	assert.Equal(t, 5.3, v)
}

func TestFloatTree_NotifyOnlyOnceWhenChange(t *testing.T) {
	f := NewFloatTree()
	triggered := 0
	f.AddListener(NewDataListener(func() {
		triggered++
	}))
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 55, "2": 77})
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.SetValue("1", 5)
	assert.Zero(t, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1", "2"}}, map[string]float64{"1": 101, "2": 98})
	assert.Zero(t, triggered)

	triggered = 0
	f.Append("1", "3", 88)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Prepend("", "4", 23)
	assert.Equal(t, 1, triggered)

	triggered = 0
	f.Set(map[string][]string{"": {"1"}}, map[string]float64{"1": 32})
	assert.Equal(t, 1, triggered)
}
