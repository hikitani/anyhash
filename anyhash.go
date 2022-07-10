package anyhash

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/hikitani/anyhash/internal"
)

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

type AnyHasher[T any] struct {
	ptrAndSizeGetters []ptrAndSizeGetter
	seed              uint
}

func (h *AnyHasher[T]) GetHash(v T) uint {
	p := noescape(unsafe.Pointer(&v))
	var seed = uintptr(h.seed)
	for _, getter := range h.ptrAndSizeGetters {
		np, sz := getter.getPtrAndSize(p)
		seed = internal.MemhashFallback(np, seed, sz)
	}

	return uint(seed)
}

type hashBuilder[T any] struct {
	h       *AnyHasher[T]
	c       cycleDeclChecker
	baseTyp reflect.Type
}

func (b *hashBuilder[T]) fill(
	v reflect.Value,
	offset uintptr,
	parentTyp reflect.Type,
	ptrDepth int,

) error {
	typ := v.Type()
	b.c.typVisits[typ.String()] = notVisiting
	if _, ok := b.c.typEdge[typ.String()]; !ok {
		b.c.typEdge[typ.String()] = map[string]struct{}{}
	}

	if parentTyp != nil {
		b.c.typEdge[parentTyp.String()][typ.String()] = struct{}{}
		if b.c.isCycle(b.baseTyp.String(), true) {
			return errors.New("anyhash: found cycle declaration")
		}
	}

	var ptrAndSizeGetter ptrAndSizeGetter
	switch k := typ.Kind(); k {
	case reflect.String:
		ptrAndSizeGetter = &stringGetter{
			offset:   offset,
			ptrDepth: ptrDepth,
		}
	case reflect.Slice:
		elemSz, err := getElemSzOfSlice(typ.Elem())
		if err != nil {
			return err
		}
		ptrAndSizeGetter = &sliceGetter{
			offset:   offset,
			ptrDepth: ptrDepth,
			elemSz:   elemSz,
		}
	case reflect.Array:
		len, elemSz, err := getLenAndElemSzArray(typ)
		if err != nil {
			return err
		}
		ptrAndSizeGetter = &arrayGetter{
			offset:   offset,
			ptrDepth: ptrDepth,
			len:      len,
			elemSz:   elemSz,
		}
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if err := b.fill(v.Field(i), field.Offset+offset, typ, ptrDepth); err != nil {
				return err
			}
		}
	case reflect.Pointer:
		if typ.Elem().Kind() == reflect.Struct {
			return errors.New("pointer of struct cannot be hashed")
		}
		if err := b.fill(reflect.Indirect(reflect.New(typ.Elem())), offset, parentTyp, ptrDepth+1); err != nil {
			return err
		}
	case reflect.Chan, reflect.Invalid, reflect.Func, reflect.Map,
		reflect.Interface, reflect.UnsafePointer:
		return fmt.Errorf("type %s cannot be hashed", k.String())
	default:
		ptrAndSizeGetter = &baseTypeGetter{
			offset:   offset,
			ptrDepth: ptrDepth,
			elemSz:   typ.Size(),
		}
	}
	if ptrAndSizeGetter != nil {
		b.h.ptrAndSizeGetters = append(b.h.ptrAndSizeGetters, ptrAndSizeGetter)
	}
	return nil
}

func elemHasPointers(elemTyp reflect.Type) bool {
	switch k := elemTyp.Kind(); k {
	case reflect.Invalid, reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Slice, reflect.String:
		return true
	}
	return false
}

func structHasPointers(structType reflect.Type) bool {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if elemHasPointers(field.Type) {
			return true
		}

		switch k := field.Type.Kind(); k {
		case reflect.Struct:
			if structHasPointers(field.Type) {
				return true
			}
		}
	}

	return false
}

func checkElemTypeOfArrayOrSlice(elemTyp reflect.Type) error {
	invalidElemErr := errors.New("element of array or slice must be basic type (bool, int, float, complex, not pointer) or struct without pointers")

	if elemHasPointers(elemTyp) {
		return invalidElemErr
	}
	switch k := elemTyp.Kind(); k {
	case reflect.Struct:
		if structHasPointers(elemTyp) {
			return invalidElemErr
		}
	}

	return nil
}

func getLenAndElemSzArray(arrTyp reflect.Type) (int, uintptr, error) {
	if k := arrTyp.Kind(); k != reflect.Array {
		return 0, 0, fmt.Errorf("expected array type, got %s", k.String())
	}

	elemTyp := arrTyp.Elem()
	len := arrTyp.Len()
	for elemTyp.Kind() == reflect.Array {
		len *= elemTyp.Len()
		elemTyp = elemTyp.Elem()
	}

	return len, elemTyp.Size(), checkElemTypeOfArrayOrSlice(elemTyp)
}

func getElemSzOfSlice(elemTyp reflect.Type) (int, error) {
	if err := checkElemTypeOfArrayOrSlice(elemTyp); err != nil {
		return 0, err
	}

	return int(elemTyp.Size()), nil
}

func NewAnyHasher[T any](seed uint) (*AnyHasher[T], error) {
	var v T
	val := reflect.ValueOf(v)
	h := &AnyHasher[T]{
		ptrAndSizeGetters: []ptrAndSizeGetter{},
		seed:              seed,
	}

	c := cycleDeclChecker{
		typEdge:   map[string]map[string]struct{}{},
		typVisits: map[string]visitStatus{},
	}

	b := hashBuilder[T]{
		c:       c,
		h:       h,
		baseTyp: val.Type(),
	}
	err := b.fill(val, 0, nil, 0)
	if err != nil {
		return nil, err
	}
	return h, nil
}
