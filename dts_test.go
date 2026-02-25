package dts

import (
	"bytes"
	"errors"
	"testing"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/dts-stream-go/testutil"
	"github.com/mus-format/mus-stream-go/testutil/mock"
	asserterror "github.com/ymz-ncnk/assert/error"
)

func TestDTS(t *testing.T) {
	t.Run("Marshal, Unmarshal, Size, Skip methods should work succeed",
		func(t *testing.T) {
			var (
				foo    = testutil.Foo{Num: 11, Str: "hello world"}
				fooDTS = New[testutil.Foo](testutil.FooDTM, testutil.FooSer)
				size   = fooDTS.Size(foo)
				buf    = bytes.NewBuffer(make([]byte, 0, size))
			)
			n, err := fooDTS.Marshal(foo, buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size)

			afoo, n, err := fooDTS.Unmarshal(buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size)
			asserterror.EqualDeep(t, afoo, foo)

			buf.Reset()

			fooDTS.Marshal(foo, buf)
			n, err = fooDTS.Skip(buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size)
		})

	t.Run("Marshal, UnmarshalDTM, UnmarshalData, Size, SkipDTM, SkipData methods should succeed",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				foo        = testutil.Foo{Num: 11, Str: "hello world"}
				fooDTS     = New[testutil.Foo](testutil.FooDTM, testutil.FooSer)
				size       = fooDTS.Size(foo)
				buf        = bytes.NewBuffer(make([]byte, 0, size))
			)
			n, err := fooDTS.Marshal(foo, buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size)

			dtm, n, err := DTMSer.Unmarshal(buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, wantDTSize)
			asserterror.Equal(t, dtm, testutil.FooDTM)

			afoo, n, err := fooDTS.UnmarshalData(buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size-wantDTSize)
			asserterror.EqualDeep(t, afoo, foo)

			buf.Reset()

			fooDTS.Marshal(foo, buf)
			_, err = DTMSer.Skip(buf)
			asserterror.EqualError(t, err, nil)

			n, err = fooDTS.SkipData(buf)
			asserterror.EqualError(t, err, nil)
			asserterror.Equal(t, n, size-wantDTSize)
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			wantDTM = testutil.FooDTM

			fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
		)

		dtm := fooDTS.DTM()
		asserterror.Equal(t, dtm, wantDTM)
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testutil.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testutil.FooDTM, actualDTM)
				wantFoo    = testutil.Foo{}

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(actualDTM)
						return
					},
				)
				fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(r)
			asserterror.EqualError(t, err, wantErr)
			asserterror.EqualDeep(t, foo, wantFoo)
			asserterror.Equal(t, n, wantDTSize)
		})

	t.Run("Skip should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				actualDTM = testutil.FooDTM + 3

				wantDTSize = 1
				wantErr    = com.NewWrongDTMError(testutil.FooDTM, actualDTM)

				r = mock.NewReader().RegisterReadByte(
					func() (b byte, err error) {
						b = byte(actualDTM)
						return
					},
				)
				fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
			)

			n, err := fooDTS.Skip(r)
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, n, wantDTSize)
		})

	t.Run("If MarshalDTM fails with an error, Marshal should return it",
		func(t *testing.T) {
			var (
				wantErr = errors.New("write byte error")

				w = mock.NewWriter().RegisterWriteByte(func(c byte) error {
					return wantErr
				})
				fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
			)
			_, err := fooDTS.Marshal(testutil.Foo{}, w)
			asserterror.EqualError(t, err, wantErr)
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
				fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
			)
			foo, n, err := fooDTS.Unmarshal(r)
			asserterror.EqualError(t, err, wantErr)
			asserterror.EqualDeep(t, foo, testutil.Foo{})
			asserterror.Equal(t, n, 0)
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
				fooDTS = New[testutil.Foo](testutil.FooDTM, nil)
			)

			n, err := fooDTS.Skip(r)
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, n, 0)
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
			asserterror.EqualError(t, err, wantErr)
			asserterror.Equal(t, dtm, com.DTM(0))
			asserterror.Equal(t, n, 0)
		})
}
