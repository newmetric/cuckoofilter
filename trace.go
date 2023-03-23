package cuckoo

type Trace struct {
	filter *Filter

	records   []record
	bucketPow uint
}

const (
	InsertOp = iota
	DeleteOp
)

type record struct {
	fp fingerprint
	i1 uint
	op int8
}

func (f *Filter) NewTrace() *Trace {
	return &Trace{
		filter: f,

		records:   make([]record, 0, f.count),
		bucketPow: f.bucketPow,
	}
}

func (t *Trace) Length() uint64 {
	return uint64(len(t.records))
}

func (t *Trace) GetRecords() []record {
	return t.records
}

func (t *Trace) Set(record record) {
	t.records = append(t.records, record)
}

func (t *Trace) Add(data []byte) {
	i1, fp := getIndexAndFingerprint(data, t.bucketPow)
	t.Set(record{fp: fp, i1: i1, op: InsertOp})
}

func (t *Trace) AddTS(entry []byte) {
	t.filter.mtx.Lock()
	defer t.filter.mtx.Unlock()
	t.Add(entry)
}

func (t *Trace) Delete(data []byte) {
	i1, fp := getIndexAndFingerprint(data, t.bucketPow)
	t.Set(record{fp: fp, i1: i1, op: DeleteOp})
}

func (t *Trace) Sync() {
	for _, record := range t.records {
		fp := record.fp
		i1 := record.i1

		switch record.op {
		case InsertOp:
			t.filter.insert(fp, i1)
		case DeleteOp:
			t.filter.delete(fp, i1)
		}
	}
}

func (t *Trace) Reset() {
	t.records = t.records[:0]
}
