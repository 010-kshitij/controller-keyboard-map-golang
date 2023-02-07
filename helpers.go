package main

func getKeysOnIndexes(binary string, mapped map[int]string) []string {
	indexes := make([]string, 0) 
	for index, bit := range binary {
		if bit == '1' {
			indexes = append(indexes, mapped[7-index])
		}
	}
	return indexes
}
