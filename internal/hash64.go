package internal

import (
	"math/bits"
	"unsafe"
)

const (
	m1 = 0xa0761d6478bd642f
	m2 = 0xe7037ed1a0b428db
	m3 = 0x8ebc6af09c88c6e3
	m4 = 0x589965cc75374cc3
	m5 = 0x1d8e4e27c47d124f
)

func memhashFallback64(p unsafe.Pointer, seed, s uintptr) uintptr {
	var a, b uintptr
	seed ^= m1
	switch {
	case s == 0:
		return seed
	case s < 4:
		a = uintptr(*(*byte)(p))
		a |= uintptr(*(*byte)(unsafe.Add(p, s>>1))) << 8
		a |= uintptr(*(*byte)(unsafe.Add(p, s-1))) << 16
	case s == 4:
		a = r4(p)
		b = a
	case s < 8:
		a = r4(p)
		b = r4(unsafe.Add(p, s-4))
	case s == 8:
		a = r8(p)
		b = a
	case s <= 16:
		a = r8(p)
		b = r8(unsafe.Add(p, s-8))
	default:
		l := s
		if l > 48 {
			seed1 := seed
			seed2 := seed
			for ; l > 48; l -= 48 {
				seed = mix(r8(p)^m2, r8(unsafe.Add(p, 8))^seed)
				seed1 = mix(r8(unsafe.Add(p, 16))^m3, r8(unsafe.Add(p, 24))^seed1)
				seed2 = mix(r8(unsafe.Add(p, 32))^m4, r8(unsafe.Add(p, 40))^seed2)
				p = unsafe.Add(p, 48)
			}
			seed ^= seed1 ^ seed2
		}
		for ; l > 16; l -= 16 {
			seed = mix(r8(p)^m2, r8(unsafe.Add(p, 8))^seed)
			p = unsafe.Add(p, 16)
		}
		a = r8(unsafe.Add(p, l-16))
		b = r8(unsafe.Add(p, l-8))
	}

	return mix(m5^s, mix(a^m2, b^seed))
}

func mix(a, b uintptr) uintptr {
	hi, lo := bits.Mul64(uint64(a), uint64(b))
	return uintptr(hi ^ lo)
}

func r4(p unsafe.Pointer) uintptr {
	q := (*[4]byte)(p)
	return uintptr(Endian.Uint32(q[:]))
}

func r8(p unsafe.Pointer) uintptr {
	q := (*[8]byte)(p)
	return uintptr(Endian.Uint64(q[:]))
}
