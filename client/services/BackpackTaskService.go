package services

import (
	"client/backpackTaskGRPC"
	"fmt"
)

func SolveBackPackTask(part *backpackTaskGRPC.TaskPart) float64 {
	backpackCapacity := part.BackpackCapacity
	fmt.Println("backpackCapacity", backpackCapacity)
	itemsCount := len(part.Items)
	F := *new([][]backpackTaskGRPC.Item)
	for i := 0; i < itemsCount+1; i++ {
		var a []backpackTaskGRPC.Item
		for j := 0; j < int(backpackCapacity)+1; j++ {
			a = append(a, backpackTaskGRPC.Item{})
		}
		F = append(F, a)
	}

	for i := 1; i < itemsCount+1; i++ {
		for k := 1; uint32(k) < backpackCapacity+1; k++ {
			if uint32(k) >= part.Items[i-1].Weight {
				F[i][k].Price = MaxUint32(F[i-1][k].Price, F[i-1][k-int(part.Items[i-1].Weight)].Price+part.Items[i-1].Price)
				F[i][k].Weight = part.Items[i-1].Weight
				F[i][k].Id = part.Items[i-1].Id
			} else {
				F[i][k].Price = F[i-1][k].Price
				F[i][k].Weight = F[i-1][k].Weight
				F[i][k].Id = F[i-1][k].Id
			}
		}
	}

	var ans []backpackTaskGRPC.Item
	k := int(backpackCapacity) - 1
	for i := itemsCount; i > 0; i-- {
		if F[i][k].Price != F[i-1][k].Price {
			ans = append(ans, *part.Items[i-1])
			k -= int(part.Items[i-1].Weight)
		}
	}

	check := uint32(0)
	for _, j := range ans {
		check += j.Weight
		fmt.Println(j.Id, j.Weight, j.Price)
	}
	boba := uint32(0)
	aboba := uint32(0)
	for _, i := range part.Items {
		boba += i.Weight
		aboba += i.Price
	}
	fmt.Println("Общий вес предметов:", boba)
	fmt.Println("Общая ценность предметов:", aboba)
	fmt.Println(backpackCapacity, check)
	return 0
}
