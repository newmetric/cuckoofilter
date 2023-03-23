package cuckoo_test

import (
	"testing"

	cuckoo "github.com/newmetric/cuckoofilter"
	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	filter := cuckoo.NewFilter(1_000)
	trace := filter.NewTrace()

	trace.Add([]byte("foo"))
	trace.Add([]byte("bar"))

	trace.Sync()

	assert.Equal(t, uint64(0x2), trace.Length())

	assert.True(t, filter.Lookup([]byte("foo")))
	assert.True(t, filter.Lookup([]byte("bar")))
	assert.False(t, filter.Lookup([]byte("baz")))

	assert.True(t, filter.Delete([]byte("foo")))
	assert.False(t, filter.Delete([]byte("foo")))

	trace.Reset()
	assert.Equal(t, uint64(0x0), trace.Length())
}
