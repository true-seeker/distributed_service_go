package services

import (
	"client/backpackTaskGRPC"
	"fmt"
)

func SolveBackPackTask(task *backpackTaskGRPC.Task) float64 {
	backpackCapacity := task.BackpackCapacity
	itemsCount := len(task.Items)
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
			if uint32(k) >= task.Items[i-1].Weight {
				F[i][k].Price = MaxUint32(F[i-1][k].Price, F[i-1][k-int(task.Items[i-1].Weight)].Price+task.Items[i-1].Price)
				F[i][k].Weight = task.Items[i-1].Weight
				F[i][k].Id = task.Items[i-1].Id
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
			ans = append(ans, *task.Items[i-1])
			k -= int(task.Items[i-1].Weight)
		}
	}

	check := uint32(0)
	for _, j := range ans {
		check += j.Weight
	}
	totalWeight := uint32(0)
	totalPrice := uint32(0)
	ansWeight := uint32(0)
	ansPrice := uint32(0)
	for _, i := range task.Items {
		totalWeight += i.Weight
		totalPrice += i.Price
	}
	fmt.Println("Вместимость рюкзака:", backpackCapacity)
	fmt.Println("Предметы: ", task.Items)
	fmt.Println("Общий вес предметов:", totalWeight)
	fmt.Println("Общая ценность предметов:", totalPrice)
	fmt.Print("\nОтвет: ")
	for _, i := range ans {
		fmt.Printf("id:%d weight:%d price:%d, ", i.Id, i.Weight, i.Price)
		ansWeight += i.Weight
		ansPrice += i.Price
	}
	fmt.Println()
	fmt.Println("Вес предметов ответа:", ansWeight)
	fmt.Println("Ценность предметов ответа:", ansPrice)
	return 0
}
