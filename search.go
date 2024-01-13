package main

import (
	"errors"
	"math/big"
)

var ErrNotFound = errors.New("target not found")

type IpItem struct {
	start   *big.Int
	end     *big.Int
	country string
}

func Binary(array []IpItem, target *big.Int, lowIndex int, highIndex int) (int, error) {
	if highIndex < lowIndex || len(array) == 0 {
		return -1, ErrNotFound
	}
	mid := int(lowIndex + (highIndex-lowIndex)/2)
	if array[mid].start.Cmp(target) == 1 {
		return Binary(array, target, lowIndex, mid-1)
	} else if array[mid].end.Cmp(target) == -1 {
		return Binary(array, target, mid+1, highIndex)
	} else {
		return mid, nil
	}
}
