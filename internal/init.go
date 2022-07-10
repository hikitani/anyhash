package internal

import (
	"encoding/binary"
	"unsafe"
)

const is64Bit = uint64(^uintptr(0)) == ^uint64(0)

var (
	Endian          binary.ByteOrder
	MemhashFallback func(p unsafe.Pointer, seed, s uintptr) uintptr
)

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		Endian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		Endian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}

	if is64Bit {
		MemhashFallback = memhashFallback64
	} else {
		MemhashFallback = memhashFallback32
	}
}
