package util

func Copy2(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
