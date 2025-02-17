package cuckoo

type Fingerprint byte

type Bucket [BucketSize]Fingerprint

const (
	nullFp     = 0
	BucketSize = 4
)

func ToBytes(buckets []Bucket) [][BucketSize]byte {
	bucket := make([][BucketSize]byte, len(buckets))

	for i, b := range buckets {
		for j, f := range b {
			bucket[i][j] = byte(f)
		}
	}
	return bucket
}

func (b *Bucket) insert(fp Fingerprint) bool {
	for i, tfp := range b {
		if tfp == nullFp {
			b[i] = fp
			return true
		}
	}
	return false
}

func (b *Bucket) delete(fp Fingerprint) bool {
	for i, tfp := range b {
		if tfp == fp {
			b[i] = nullFp
			return true
		}
	}
	return false
}

func (b *Bucket) getFingerprintIndex(fp Fingerprint) int {
	for i, tfp := range b {
		if tfp == fp {
			return i
		}
	}
	return -1
}

func (b *Bucket) reset() {
	for i := range b {
		b[i] = nullFp
	}
}
