package random

import (
	"math/rand"
	"time"
)

func PickFromIntSlice(list []int) int {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	i := r.Intn(len(list))
	return list[i]
}

func PickFromStringSlice(list []string) string {
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	i := r.Intn(len(list))
	return list[i]
}
