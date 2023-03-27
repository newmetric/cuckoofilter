package cuckoo

type Trace struct {
	filter *Filter

	records   []Record
	bucketPow uint
}

const (
	InsertOp = iota
	DeleteOp
)

type Record struct {
	Fp byte
	I1 uint
	Op int8
}

func (f *Filter) NewTrace() *Trace {
	return &Trace{
		filter: f,

		records:   make([]Record, 0, f.count),
		bucketPow: f.bucketPow,
	}
}

func (t *Trace) Length() uint64 {
	return uint64(len(t.records))
}

func (t *Trace) GetRecords() []Record {
	return t.records
}

func (t *Trace) Set(record Record) {
	t.records = append(t.records, record)
}

func (t *Trace) Add(data []byte) {
	i1, fp := GetIndexAndFingerprint(data, t.bucketPow)
	t.Set(Record{Fp: byte(fp), I1: i1, Op: InsertOp})
}

func (t *Trace) AddTS(entry []byte) {
	t.filter.mtx.Lock()
	defer t.filter.mtx.Unlock()
	t.Add(entry)
}

func (t *Trace) Delete(data []byte) {
	i1, fp := GetIndexAndFingerprint(data, t.bucketPow)
	t.Set(Record{Fp: byte(fp), I1: i1, Op: DeleteOp})
}

func (t *Trace) DeleteTS(entry []byte) {
	t.filter.mtx.Lock()
	defer t.filter.mtx.Unlock()
	t.Delete(entry)
}

func (t *Trace) Sync() {
	for _, record := range t.records {
		fp := record.Fp
		i1 := record.I1

		switch record.Op {
		case InsertOp:
			t.filter.insert(Fingerprint(fp), i1)
		case DeleteOp:
			t.filter.delete(Fingerprint(fp), i1)
		}
	}
}

func (t *Trace) SyncTS() {
	t.filter.mtx.Lock()
	defer t.filter.mtx.Unlock()
	t.Sync()
}

func (t *Trace) Reset() {
	t.records = t.records[:0]
}
