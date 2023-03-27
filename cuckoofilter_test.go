package cuckoo

import (
	"bufio"
	"crypto/rand"
	"io"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertion(t *testing.T) {
	if runtime.GOOS != "darwin" {
		t.SkipNow()
	}

	cf := NewFilter(1_000_000)
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	var lineCount uint
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if cf.InsertUnique(s) {
			lineCount++
		}
		values = append(values, s)
	}

	count := cf.Count()
	if count != lineCount {
		t.Errorf("Expected count = %d, instead count = %d", lineCount, count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.Count()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}

func TestGetBuckets(t *testing.T) {
	cf := NewFilter(8)
	cf.buckets = []Bucket{
		[4]fingerprint{1, 2, 3, 4},
		[4]fingerprint{5, 6, 7, 8},
	}
	cf.count = 8
	buckets := cf.GetBuckets()

	for i, b := range buckets {
		for j, f := range b {
			assert.Equal(t, f, cf.buckets[i][j])
		}
	}
}

func TestReplaceBuckets(t *testing.T) {
	cf := NewFilter(8)
	cf.Insert([]byte{1, 2, 3, 4})
	cf.Insert([]byte{5, 6, 7, 8})

	assert.True(t, cf.Lookup([]byte{1, 2, 3, 4}))

	othercf := NewFilter(8)
	othercf.Insert([]byte{9, 10, 11, 12})
	othercf.Insert([]byte{13, 14, 15, 16})

	assert.True(t, othercf.Lookup([]byte{9, 10, 11, 12}))

	cf.ReplaceBuckets(othercf.GetBuckets())

	assert.True(t, cf.Lookup([]byte{9, 10, 11, 12}))
	assert.False(t, cf.Lookup([]byte{1, 2, 3, 4}))
}

func TestEncodeDecode(t *testing.T) {
	cf := NewFilter(8)
	cf.buckets = []Bucket{
		[4]fingerprint{1, 2, 3, 4},
		[4]fingerprint{5, 6, 7, 8},
	}
	cf.count = 8
	bytes := cf.Encode()
	ncf, err := Decode(bytes)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(cf, ncf) {
		t.Errorf("Expected %v, got %v", cf, ncf)
	}
}

func TestDecode(t *testing.T) {
	ncf, err := Decode([]byte(""))
	if err == nil {
		t.Errorf("Expected err, got nil")
	}
	if ncf != nil {
		t.Errorf("Expected nil, got %v", ncf)
	}
}

func BenchmarkFilter_Reset(b *testing.B) {
	const cap = 10_000
	filter := NewFilter(cap)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Reset()
	}
}

func BenchmarkFilter_Insert(b *testing.B) {
	const cap = 10_000
	filter := NewFilter(cap)

	b.ResetTimer()

	var hash [32]byte
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}
}

func BenchmarkFilter_Lookup(b *testing.B) {
	const cap = 10_000
	filter := NewFilter(cap)

	var hash [32]byte
	for i := 0; i < 10_000; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Insert(hash[:])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		io.ReadFull(rand.Reader, hash[:])
		filter.Lookup(hash[:])
	}
}
