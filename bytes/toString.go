package libraryBytes

func BytesToBitString(data []byte) string {
	result := ""
	for _, b := range data {
		for i := 7; i >= 0; i-- {
			if (b>>i)&1 == 1 {
				result += "1"
			} else {
				result += "0"
			}
		}
	}
	return result
}
