// Package dtss provides DTM (Deterministic Type Marker) support for
// mus-stream-go serializer. It wraps a mus-stream-go serializer together with
// a DTM value, enabling serialization of data with an embedded type marker.
package dtss

import (
	com "github.com/mus-format/common-go"
	muss "github.com/mus-format/mus-stream-go"
)

// New creates a new DTS.
func New[T any](dtm com.DTM, ser muss.Serializer[T],
) DTS[T] {
	return DTS[T]{dtm, ser}
}

// DTS implements the mus.Serializer interface and provides DTM support for the
// mus-stream-go serializer. It helps to serializer DTM + data.
type DTS[T any] struct {
	dtm com.DTM
	ser muss.Serializer[T]
}

// DTM returns the initialization value.
func (d DTS[T]) DTM() com.DTM {
	return d.dtm
}

// Marshal marshals DTM + data.
func (d DTS[T]) Marshal(t T, w muss.Writer) (n int, err error) {
	n, err = DTMSer.Marshal(d.dtm, w)
	if err != nil {
		return
	}
	var n1 int
	n1, err = d.ser.Marshal(t, w)
	n += n1
	return
}

// Unmarshal unmarshals DTM + data.
//
// Returns ErrWrongDTM if the unmarshalled DTM differs from the dts.DTM().
func (d DTS[T]) Unmarshal(r muss.Reader) (t T, n int, err error) {
	dtm, n, err := DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	if dtm != d.dtm {
		err = ErrWrongDTM
		return
	}
	var n1 int
	t, n1, err = d.UnmarshalData(r)
	n += n1
	return
}

// Size calculates the size of the DTM + data.
func (d DTS[T]) Size(t T) (size int) {
	size = DTMSer.Size(d.dtm)
	return size + d.ser.Size(t)
}

// Skip skips DTM + data.
//
// Returns ErrWrongDTM if the unmarshalled DTM differs from the dts.DTM().
func (d DTS[T]) Skip(r muss.Reader) (n int, err error) {
	dtm, n, err := DTMSer.Unmarshal(r)
	if err != nil {
		return
	}
	if dtm != d.dtm {
		err = ErrWrongDTM
		return
	}
	n1, err := d.SkipData(r)
	n += n1
	return
}

// UnmarshalData unmarshals only data.
func (d DTS[T]) UnmarshalData(r muss.Reader) (t T, n int, err error) {
	return d.ser.Unmarshal(r)
}

// SkipData skips only data.
func (d DTS[T]) SkipData(r muss.Reader) (n int, err error) {
	return d.ser.Skip(r)
}
