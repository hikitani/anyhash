package anyhash

import (
	"fmt"
	"testing"

	"github.com/hikitani/anyhash/internal"
)

type testFoo struct {
	str  string
	b    byte
	i    int
	i16  int16
	ui32 uint32
	bs   []byte
}

func TestBaseTypes(t *testing.T) {
	v := int64(-2)
	b := make([]byte, 8)
	internal.Endian.PutUint64(b, uint64(v))
	fmt.Print(b)
	t.Run("Bool", testType(true, []byte{1}))
	t.Run("Int8", testType(int8(-1), []byte{255}))
	t.Run("Int16", testType(int16(25252), []byte{164, 98}))
	t.Run("Int32", testType(int32(252525252), []byte{196, 58, 13, 15}))
	t.Run("Int64", testType(int64(2525252525252525252), []byte{196, 188, 213, 215, 202, 127, 11, 35}))
	t.Run("UInt8", testType(uint8(255), []byte{255}))
	t.Run("UInt16", testType(uint16(25252), []byte{164, 98}))
	t.Run("UInt32", testType(uint32(252525252), []byte{196, 58, 13, 15}))
	t.Run("UInt64", testType(uint64(2525252525252525252), []byte{196, 188, 213, 215, 202, 127, 11, 35}))
	t.Run("Float32", testType(float32(232323232323.1212121212), []byte{45, 94, 88, 82}))
	t.Run("Float64", testType(float64(-48592615485926154815261548.9584516295845162), []byte{126, 197, 94, 123, 241, 24, 68, 197}))
	t.Run("Complex64", testType(complex64(232323232323.1212121212+232323232323.1212121212i), []byte{45, 94, 88, 82, 45, 94, 88, 82}))
	t.Run("Complex128", testType(complex128(-48592615485926154815261548.9584516295845162+-48592615485926154815261548.9584516295845162i), []byte{126, 197, 94, 123, 241, 24, 68, 197, 126, 197, 94, 123, 241, 24, 68, 197}))
	t.Run("String", testType("Hello, world!", []byte("Hello, world!")))
	t.Run("Array[int16]", testType([4]int16{-1, -2, 11111}, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Array[int16][int16]", testType([2][2]int16{{-1, -2}, {11111}}, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Slice[int16]", testType([]int16{-1, -2, 11111}, []byte{255, 255, 254, 255, 103, 43}))
}

func TestStructs(t *testing.T) {
	t.Run("simple", testType(struct {
		a int8
		b int16
	}{
		a: -1,
		b: 25252,
	}, []byte{255}, []byte{164, 98}))

	s := "string"
	t.Run("nested", testType(struct {
		a int8
		b int16
		c struct {
			s string
			a struct {
				b []byte
			}
			sp *string
		}
	}{
		a: -1,
		b: 25252,
		c: struct {
			s  string
			a  struct{ b []byte }
			sp *string
		}{
			s: "string 1",
			a: struct{ b []byte }{
				b: []byte("bytes"),
			},
			sp: &s,
		},
	}, []byte{255}, []byte{164, 98}, []byte("string 1"), []byte("bytes"), []byte(s)))
}

func testType[T any](v T, inBytes ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		hv, err := NewAnyHasher[T](0)
		if err != nil {
			t.Fatalf("expected nil err, got %s", err)
		}

		got := hv.GetHash(v)

		want := uint(0)
		for _, b := range inBytes {
			hb, err := NewAnyHasher[[]byte](uint(want))
			if err != nil {
				t.Fatalf("expected nil err, got %s", err)
			}
			want = hb.GetHash(b)
		}
		if got != want {
			t.Fatalf("got %d, want %d", got, want)
		}
	}
}
