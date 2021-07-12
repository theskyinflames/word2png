package tooling

// SplitByte splits a given byte in two other bytes,
// one with its high bits, and another with its low bits
func SplitByte(b byte) (high byte, low byte) {
	low = b & 0x0F
	high = b >> 4
	return
}

// JoinByte takes high and low parts and join them in a byte
func JoinByte(high, low byte) (originalByte byte) {
	return high<<4 | low
}

// Taken from https://stackoverflow.com/questions/52811744/extract-bits-into-a-int-slice-from-byte-slice
func bytes2bits(data []byte) []int8 {
	r := make([]int8, len(data)*8)
	for i, b := range data {
		for j := 0; j < 8; j++ {
			r[i*8+j] = int8(b >> uint(7-j) & 0x01)
		}
	}
	return r
}
