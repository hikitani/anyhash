package internal

import "unsafe"

func memhashFallback32(p unsafe.Pointer, seed, s uintptr) uintptr {
	a, b := mix32(uint32(seed^(s>>32)), uint32(s))
	if s == 0 {
		return uintptr(a ^ b)
	}
	for ; s > 8; s -= 8 {
		a ^= uint32(r4(p))
		b ^= uint32(r4(unsafe.Add(p, 4)))
		a, b = mix32(a, b)
		p = unsafe.Add(p, 8)
	}
	if s >= 4 {
		a ^= uint32(r4(p))
		b ^= uint32(r4(unsafe.Add(p, s-4)))
	} else {
		t := uint32(*(*byte)(p))
		t |= uint32(*(*byte)(unsafe.Add(p, s>>1))) << 8
		t |= uint32(*(*byte)(unsafe.Add(p, s-1))) << 16
		b ^= t
	}
	a, b = mix32(a, b)
	a, b = mix32(a, b)
	return uintptr(a ^ b)
}

func mix32(a, b uint32) (uint32, uint32) {
	c := uint64(a^0x53c5ca59) * uint64(b^0x74743c1b)
	return uint32(c), uint32(c >> 32)
}
