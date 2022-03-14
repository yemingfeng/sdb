package collection

type Row struct {
	Key     []byte
	Id      []byte
	Value   []byte
	Indexes []Index
}

type Index struct {
	Name  []byte
	Value []byte
}
