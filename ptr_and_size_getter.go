package anyhash

import (
	"fmt"
	"reflect"
	"unsafe"
)

func indirect(p unsafe.Pointer, depth int) unsafe.Pointer {
	switch depth {
	case 0:
		return p
	case 1:
		return *(*unsafe.Pointer)(p)
	}

	return indirect(*(*unsafe.Pointer)(p), depth-1)
}

type ptrAndSizeGetter interface {
	getPtrAndSize(p unsafe.Pointer) (unsafe.Pointer, uintptr)
}

type baseTypeGetter struct {
	offset   uintptr
	ptrDepth int
	elemSz   uintptr
}

func (b *baseTypeGetter) getPtrAndSize(p unsafe.Pointer) (unsafe.Pointer, uintptr) {
	np := indirect(unsafe.Pointer(uintptr(p)+b.offset), b.ptrDepth)
	fmt.Println(*(*bool)(np))
	return np, b.elemSz
}

type stringGetter struct {
	offset   uintptr
	ptrDepth int
}

func (s *stringGetter) getPtrAndSize(p unsafe.Pointer) (unsafe.Pointer, uintptr) {
	np := indirect(unsafe.Pointer(uintptr(p)+s.offset), s.ptrDepth)
	sh := (*reflect.StringHeader)(np)
	return unsafe.Pointer(sh.Data), uintptr(sh.Len)
}

type sliceGetter struct {
	offset   uintptr
	ptrDepth int
	elemSz   int
}

func (s *sliceGetter) getPtrAndSize(p unsafe.Pointer) (unsafe.Pointer, uintptr) {
	np := indirect(unsafe.Pointer(uintptr(p)+s.offset), s.ptrDepth)
	sh := (*reflect.SliceHeader)(np)
	return unsafe.Pointer(sh.Data), uintptr(sh.Len * s.elemSz)
}

type arrayGetter struct {
	offset   uintptr
	ptrDepth int
	len      int
	elemSz   uintptr
}

func (a *arrayGetter) getPtrAndSize(p unsafe.Pointer) (unsafe.Pointer, uintptr) {
	np := indirect(unsafe.Pointer(uintptr(p)+a.offset), a.ptrDepth)
	return np, uintptr(a.len) * a.elemSz
}
