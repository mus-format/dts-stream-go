package dtss

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/mus-format/dts-stream-go/testdata"
	muss_mock "github.com/mus-format/mus-stream-go/testdata/mock"
)

func TestDTS(t *testing.T) {
	t.Run("Marshal, Unmarshal, Size, Skip methods should work correctly",
		func(t *testing.T) {
			var (
				foo    = testdata.Foo{Num: 11, Str: "hello world"}
				fooDTS = New[testdata.Foo](testdata.FooDTM, testdata.FooSer)
				size   = fooDTS.Size(foo)
				buf    = bytes.NewBuffer(make([]byte, 0, size))
			)
			n, err := fooDTS.Marshal(foo, buf)
			if err != nil {
				t.Fatal(err)
			}
			if n != size {
				t.Fatalf("unexpected n, want '%v' actual '%v'", size, n)
			}
			afoo, n, err := fooDTS.Unmarshal(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != size {
				t.Errorf("unexpected n, want '%v' actual '%v'", size, n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			buf.Reset()
			fooDTS.Marshal(foo, buf)
			n, err = fooDTS.Skip(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != size {
				t.Errorf("unexpected n, want '%v' actual '%v'", size, n)
			}
		})

	t.Run("Marshal, UnmarshalDTM, UnmarshalData, Size, SkipDTM, SkipData methods should work correctly",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				foo        = testdata.Foo{Num: 11, Str: "hello world"}
				fooDTS     = New[testdata.Foo](testdata.FooDTM, testdata.FooSer)
				size       = fooDTS.Size(foo)
				buf        = bytes.NewBuffer(make([]byte, 0, size))
			)
			n, err := fooDTS.Marshal(foo, buf)
			if err != nil {
				t.Fatal(err)
			}
			if n != size {
				t.Fatalf("unexpected n, want '%v' actual '%v'", size, n)
			}
			dtm, n, err := DTMSer.Unmarshal(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", 1, n)
			}
			if dtm != testdata.FooDTM {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", testdata.FooDTM, dtm)
			}
			afoo, n, err := fooDTS.UnmarshalData(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != size-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", size, n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			buf.Reset()
			fooDTS.Marshal(foo, buf)
			_, err = DTMSer.Skip(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			n, err = fooDTS.SkipData(buf)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != size-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", size, n)
			}
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			dtm    = fooDTS.DTM()
		)
		if dtm != testdata.FooDTM {
			t.Errorf("unexpected dtm, want '%v' actual '%v'", testdata.FooDTM, dtm)
		}
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				r          = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(testdata.FooDTM) + 3
						return
					},
				)
				fooDTS      = New[testdata.Foo](testdata.FooDTM, nil)
				foo, n, err = fooDTS.Unmarshal(r)
			)
			if err != ErrWrongDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrWrongDTM, err)
			}
			if !reflect.DeepEqual(foo, testdata.Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", wantDTSize, n)
			}
		})

	t.Run("Skip should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				r          = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(testdata.FooDTM) + 3
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
				n, err = fooDTS.Skip(r)
			)
			if err != ErrWrongDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrWrongDTM, err)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", wantDTSize, n)
			}
		})

	t.Run("If MarshalDTM fails with an error, Marshal should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("write byte error")
				w       = muss_mock.NewWriter().RegisterWriteByte(func(c byte) error {
					return wantErr
				})
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			_, err := fooDTS.Marshal(testdata.Foo{}, w)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", wantErr, err)
			}
		})

	t.Run("If UnmarshalDTM fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")
				r       = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				fooDTS      = New[testdata.Foo](testdata.FooDTM, nil)
				foo, n, err = fooDTS.Unmarshal(r)
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", io.EOF, err)
			}
			if !reflect.DeepEqual(foo, testdata.Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If UnmarshalDTM fails with an error, Skip should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")
				r       = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
				n, err = fooDTS.Skip(r)
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", io.EOF, err)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If varint.UnmarshalInt fails with an error, UnmarshalDTM should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")
				r       = muss_mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				dtm, n, err = DTMSer.Unmarshal(r)
			)
			if err != wantErr {
				t.Errorf("unexpected error, want '%v' actual '%v'", io.EOF, err)
			}
			if dtm != 0 {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", 0, dtm)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})
}
