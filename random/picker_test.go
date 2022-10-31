package random

import "testing"

func TestPickFromSlice(t *testing.T) {
	reasons := make([]string, 0)
	reasons = append(reasons,
		"Locked out",
		"Pipes broke",
		"Food poisoning",
		"Not feeling well")

	t.Log(PickFromStringSlice(reasons))

	num := make([]int, 0)
	num = append(num, 1, 2, 3)

	t.Log(PickFromIntSlice(num))
}
