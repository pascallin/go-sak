package random

import "math/rand"

func GenerateIntSlice(length int) []int {
	var list []int
	if length <= 0 {
		return list
	}

	for i := 0; i < length; i++ {
		list = append(list, i)
	}
	for i := len(list) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		list[i], list[num] = list[num], list[i]
	}
	return list
}
