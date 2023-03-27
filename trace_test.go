package cuckoo_test

import (
	"testing"

	cuckoo "github.com/newmetric/cuckoofilter"
	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	filter := cuckoo.NewFilter(1_000)

	filter.Insert([]byte("foo"))
	assert.True(t, filter.Lookup([]byte("foo")))

	trace := filter.NewTrace()

	trace.Add([]byte("bar"))
	trace.Add([]byte("baz"))
	trace.Delete([]byte("foo"))

	assert.Equal(t, uint64(3), trace.Length())

	trace.Sync()

	assert.False(t, filter.Lookup([]byte("foo")))
	assert.True(t, filter.Lookup([]byte("bar")))
	assert.True(t, filter.Lookup([]byte("baz")))

	assert.True(t, filter.Delete([]byte("bar")))
	assert.False(t, filter.Delete([]byte("bar")))

	trace.Reset()
	assert.Equal(t, uint64(0), trace.Length())
}
