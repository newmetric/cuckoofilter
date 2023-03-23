package cuckoo

type fingerprint byte

type bucket [bucketSize]fingerprint

const (
	nullFp     = 0
	bucketSize = 4
)

func getAll(buckets []bucket) [][bucketSize]byte {
	bucket := make([][bucketSize]byte, len(buckets))

	for i, b := range buckets {
		for j, f := range b {
			bucket[i][j] = byte(f)
		}
	}
	return bucket
}

func (b *bucket) insert(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == nullFp {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b *bucket) delete(fp fingerprint) bool {
	for i, tfp := range b {
		if tfp == fp {
			b[i] = nullFp
			return true
		}
	}
	return false
}

func (b *bucket) getFingerprintIndex(fp fingerprint) int {
	for i, tfp := range b {
		if tfp == fp {
			return i
		}
	}
	return -1
}

func (b *bucket) reset() {
	for i := range b {
		b[i] = nullFp
	}
}
