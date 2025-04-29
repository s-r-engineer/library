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

func HowMany(g int) (h [][]int) {
	if g <= 0 {
		return nil
	}
	for _, size := range chunkSizes {
		if g >= size {
			count := g / size
			h = append(h, []int{size, count})
			g %= size
		}
	}
	if g > 0 {
		h = append(h, []int{g, 1})
	}
	return
}
