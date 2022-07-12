package internal

import (
	"encoding/binary"
	"unsafe"
)

var endian binary.ByteOrder

func init() {
	buf := [2]byte{}
	*(*uint16)(unsafe.Pointer(&buf[0])) = uint16(0xABCD)

	switch buf {
	case [2]byte{0xCD, 0xAB}:
		endian = binary.LittleEndian
	case [2]byte{0xAB, 0xCD}:
		endian = binary.BigEndian
	default:
		panic("Could not determine native endianness.")
	}
}

func r4(p unsafe.Pointer) uintptr {
	q := (*[4]byte)(p)
	return uintptr(endian.Uint32(q[:]))
}

func r8(p unsafe.Pointer) uintptr {
	q := (*[8]byte)(p)
	return uintptr(endian.Uint64(q[:]))
}
