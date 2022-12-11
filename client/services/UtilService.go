package services

func MaxUint32(nums ...uint32) uint32 {
	max := uint32(0)

	for _, elem := range nums {
		if elem > max {
			max = elem
		}
	}
	return max
}
