package anyhash

import (
	"fmt"
	"testing"
	"unsafe"
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
	v := int16(25252)
	s := "Hello, world!"
	arr := [4]int16{-1, -2, 11111}
	slice := arr[:]
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
	t.Run("String", testType(s, []byte(s)))
	t.Run("Array[int16]", testType(arr, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Array[int16][int16]", testType([2][2]int16{{-1, -2}, {11111}}, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Slice[int16]", testType(slice, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Pointer[int16]", testType(&v, []byte{164, 98}))
	t.Run("Pointer[string]", testType(&s, []byte(s)))
	t.Run("Pointer[Array[int16]]", testType(&arr, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
	t.Run("Pointer[Slice[int16]]", testType(&slice, []byte{255, 255, 254, 255, 103, 43, 0, 0}))
}

func TestStructs(t *testing.T) {
	t.Run("Simple", testType(struct {
		a int8
		b int16
	}{
		a: -1,
		b: 25252,
	}, []byte{255}, []byte{164, 98}))

	s := "string"
	t.Run("Nested", testType(struct {
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

	t.Run("Slice[Struct]", testType([]struct {
		a int16
		b [2]byte
	}{
		{
			a: 25252,
			b: [2]byte{1, 2},
		},
		{
			a: -1,
			b: [2]byte{10, 11},
		},
	}, []byte{164, 98, 1, 2, 255, 255, 10, 11}))

	t.Run("Array[struct]", testType([2]struct {
		a int16
		b [2]byte
	}{
		{
			a: 25252,
			b: [2]byte{1, 2},
		},
		{
			a: -1,
			b: [2]byte{10, 11},
		},
	}, []byte{164, 98, 1, 2, 255, 255, 10, 11}))
}

func TestDisallowedTypes(t *testing.T) {
	t.Run("Pointer[Struct]", testDisallowedType(&struct{}{}))
	t.Run("Struct[Pointer[Struct]]", testDisallowedType(struct{ s *struct{} }{}))
	t.Run("Chan", testDisallowedType(make(chan struct{})))
	t.Run("Func", testDisallowedType(func() {}))
	t.Run("Map", testDisallowedType(map[string]struct{}{}))
	t.Run("Interface", testDisallowedType(getAny(struct{}{})))
	t.Run("UnsafePointer", testDisallowedType(unsafe.Pointer(&struct{}{})))

	t.Run("Slice[Chan]", testDisallowedType([]chan struct{}{}))
	t.Run("Slice[Func]", testDisallowedType([]func(){}))
	t.Run("Slice[Interface]", testDisallowedType([]any{}))
	t.Run("Slice[Map]", testDisallowedType([]map[string]struct{}{}))
	t.Run("Slice[Pointer]", testDisallowedType([]*struct{}{}))
	t.Run("Slice[UnsafePointer]", testDisallowedType([]unsafe.Pointer{}))
	t.Run("Slice[Slice]", testDisallowedType([][]struct{}{}))
	t.Run("Slice[String]", testDisallowedType([]string{}))
	t.Run("Slice[Struct[Chan]]", testDisallowedType([]struct{ ch chan struct{} }{}))
	t.Run("Slice[Struct[Func]]", testDisallowedType([]struct{ f func() }{}))
	t.Run("Slice[Struct[Interface]]", testDisallowedType([]struct{ a any }{}))
	t.Run("Slice[Struct[Map]]", testDisallowedType([]struct{ m map[string]struct{} }{}))
	t.Run("Slice[Struct[Pointer]]", testDisallowedType([]struct{ p *int }{}))
	t.Run("Slice[Struct[UnsafePointer]]", testDisallowedType([]struct{ p unsafe.Pointer }{}))
	t.Run("Slice[Struct[Slice]]", testDisallowedType([]struct{ s []struct{} }{}))
	t.Run("Slice[Struct[String]]", testDisallowedType([]struct{ s string }{}))

	t.Run("Array[Chan]", testDisallowedType([4]chan struct{}{}))
	t.Run("Array[Func]", testDisallowedType([4]func(){}))
	t.Run("Array[Interface]", testDisallowedType([4]any{}))
	t.Run("Array[Map]", testDisallowedType([4]map[string]struct{}{}))
	t.Run("Array[Pointer]", testDisallowedType([4]*struct{}{}))
	t.Run("Array[UnsafePointer]", testDisallowedType([4]unsafe.Pointer{}))
	t.Run("Array[Slice]", testDisallowedType([4][]struct{}{}))
	t.Run("Array[String]", testDisallowedType([4]string{}))
	t.Run("Array[Struct[Chan]]", testDisallowedType([4]struct{ ch chan struct{} }{}))
	t.Run("Array[Struct[Func]]", testDisallowedType([4]struct{ f func() }{}))
	t.Run("Array[Struct[Interface]]", testDisallowedType([4]struct{ a any }{}))
	t.Run("Array[Struct[Map]]", testDisallowedType([4]struct{ m map[string]struct{} }{}))
	t.Run("Array[Struct[Pointer]]", testDisallowedType([4]struct{ p *int }{}))
	t.Run("Array[Struct[UnsafePointer]]", testDisallowedType([4]struct{ p unsafe.Pointer }{}))
	t.Run("Array[Struct[Slice]]", testDisallowedType([4]struct{ s []struct{} }{}))
	t.Run("Array[Struct[String]]", testDisallowedType([4]struct{ s string }{}))
}

func getAny(v any) any {
	return v
}

func testType[T any](v T, inBytes ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		hv, err := New[T](0)
		if err != nil {
			t.Fatalf("expected nil err, got %s", err)
		}

		got := hv.GetHash(v)

		want := uint(0)
		for _, b := range inBytes {
			hb, err := New[[]byte](uint(want))
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

func testDisallowedType[T any](v T) func(t *testing.T) {
	return func(t *testing.T) {
		h, err := New[T](0)
		if err == nil {
			t.Fatal("err is nil")
		}

		if h != nil {
			t.Fatal("returned object is not nil")
		}
	}
}

func BenchmarkAnyHasher(b *testing.B) {
	h, err := New[[]byte](0)
	if err != nil {
		b.Fatal(err)
	}

	for i := 2; i < 17; i++ {
		l := 1 << i
		b.Run(fmt.Sprintf("%dBytes", l), func(b *testing.B) {
			bs := make([]byte, l)
			b.SetBytes(int64(l))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				h.GetHash(bs)
			}
		})
	}
}
