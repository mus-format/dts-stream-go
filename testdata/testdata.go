package testdata

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-stream-go"
	"github.com/mus-format/mus-stream-go/ord"
	"github.com/mus-format/mus-stream-go/varint"
)

const FooDTM com.DTM = 1

type Foo struct {
	Num int
	Str string
}

var FooSer = fooSer{}

type fooSer struct{}

func (s fooSer) Marshal(foo Foo, w mus.Writer) (n int, err error) {
	n, err = varint.Int.Marshal(foo.Num, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.String.Marshal(foo.Str, w)
	n += n1
	return
}

func (s fooSer) Unmarshal(r mus.Reader) (foo Foo, n int, err error) {
	foo.Num, n, err = varint.Int.Unmarshal(r)
	if err != nil {
		return
	}
	var n1 int
	foo.Str, n1, err = ord.String.Unmarshal(r)
	n += n1
	return
}

func (s fooSer) Size(foo Foo) (size int) {
	size = varint.Int.Size(foo.Num)
	return size + ord.String.Size(foo.Str)
}

func (s fooSer) Skip(r mus.Reader) (n int, err error) {
	n, err = varint.Int.Skip(r)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.String.Skip(r)
	n += n1
	return
}
