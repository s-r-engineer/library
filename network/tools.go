package libraryNetwork

var (
	K  = 1024
	M  = K * 1024
	G  = M * 1024
	hK = K / 2
	qK = hK / 2
	eK = qK / 2

	chunkSizes = []int{G, M, K, hK, qK, eK}
)

func HowMany[T uint | uint32 | uint64](g T) (h [][]int) {
	for _, size := range chunkSizes {
		if g >= T(size) {
			count := g / T(size)
			h = append(h, []int{size, int(count)})
			g %= T(size)
		}
	}
	if g > 0 {
		h = append(h, []int{int(g), 1})
	}
	return
}
