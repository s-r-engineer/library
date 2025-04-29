package libraryNetwork

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var origins = map[int][][]int{
	6013583879: {
		{G, 5},
		{M, 615},
		{hK, 1},
		{7, 1},
	},
	6297606: {
		{1048576, 6},
		{1024, 6},
		{6, 1},
	},
	1025: {
		{K, 1},
		{1, 1},
	},
	4: {
		{4, 1},
	},
	196: {
		{eK, 1},
		{68, 1},
	},
	0:  nil,
	-1: nil,
}

func TestHowMany(t *testing.T) {
	for k, v := range origins {
		result := HowMany(k)
		require.Equal(t, result, v)
		result1 := 0
		for _, vv := range v {
			result1 += vv[0] * vv[1]
		}
		if k > 0 {
			require.Equal(t, result1, k)
		}
	}
}
