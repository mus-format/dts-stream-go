package dts

import (
	"bytes"
	"errors"
	"testing"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/dts-stream-go/testdata"
	"github.com/mus-format/mus-stream-go/testdata/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
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
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size, t)

			afoo, n, err := fooDTS.Unmarshal(buf)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size, t)
			asserterror.EqualDeep(afoo, foo, t)

			buf.Reset()

			fooDTS.Marshal(foo, buf)
			n, err = fooDTS.Skip(buf)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size, t)
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
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size, t)

			dtm, n, err := DTMSer.Unmarshal(buf)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, wantDTSize, t)
			asserterror.Equal(dtm, testdata.FooDTM, t)

			afoo, n, err := fooDTS.UnmarshalData(buf)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size-wantDTSize, t)
			asserterror.EqualDeep(afoo, foo, t)

			buf.Reset()

			fooDTS.Marshal(foo, buf)
			_, err = DTMSer.Skip(buf)
			asserterror.EqualError(err, nil, t)

			n, err = fooDTS.SkipData(buf)
			asserterror.EqualError(err, nil, t)
			asserterror.Equal(n, size-wantDTSize, t)
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			wantDTM = testdata.FooDTM

			fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
		)

		dtm := fooDTS.DTM()
		asserterror.Equal(dtm, wantDTM, t)
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testdata.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testdata.FooDTM, actualDTM)
				wantFoo    = testdata.Foo{}

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(actualDTM)
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(r)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(foo, wantFoo, t)
			asserterror.Equal(n, wantDTSize, t)
		})

	t.Run("Skip should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testdata.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testdata.FooDTM, actualDTM)

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(actualDTM)
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)

			n, err := fooDTS.Skip(r)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(n, wantDTSize, t)
		})

	t.Run("If MarshalDTM fails with an error, Marshal should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("write byte error")

				w = mock.NewWriter().RegisterWriteByte(func(c byte) error {
					return wantErr
				})
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			_, err := fooDTS.Marshal(testdata.Foo{}, w)
			asserterror.EqualError(err, wantErr, t)
		})

	t.Run("If UnmarshalDTM fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(r)
			asserterror.EqualError(err, wantErr, t)
			asserterror.EqualDeep(foo, testdata.Foo{}, t)
			asserterror.Equal(n, 0, t)
		})

	t.Run("If UnmarshalDTM fails with an error, Skip should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
				fooDTS = New[testdata.Foo](testdata.FooDTM, nil)
			)

			n, err := fooDTS.Skip(r)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(n, 0, t)
		})

	t.Run("If varint.UnmarshalInt fails with an error, UnmarshalDTM should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("read byte error")

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						err = wantErr
						return
					},
				)
			)
			dtm, n, err := DTMSer.Unmarshal(r)
			asserterror.EqualError(err, wantErr, t)
			asserterror.Equal(dtm, com.DTM(0), t)
			asserterror.Equal(n, 0, t)
		})
}
