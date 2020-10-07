package ippool

func Exp2nUInt64(a int) uint64 {
	// will give (2^64)-1 for 2^64 due to memory boundary - fine for us
	var i uint64
	var b uint64
	b = 1
	for i = 0; i < (uint64(a)); i++ {
		if i < 63 {
			b = b * 2
		} else if i == 63 {
			b = (b * 2) - 1
		}
	}
	return b
}
