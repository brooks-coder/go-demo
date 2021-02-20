package main

import (
	"testing"
)

func Test_generateMatrix(t *testing.T) {
	t.Log(generateMatrix(5))
}

func generateMatrix(n int) [][]int {
	res := make([][]int, 0, n)
	for i := 0; i < n; i++ {
		res = append(res, make([]int, n))
	}
	directions := [][]int{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}
	t := (n + 1) / 2
	cur := 1
	for i := 0; i < t; i++ {
		res[i][i] = cur
		x, y := i, i
		for _, d := range directions {
			for k := 0; k < (n-1)-2*i; k++ {
				cur++
				if cur > n*n {
					return res
				}
				x += d[0]
				y += d[1]
				if x == y && x == i {
					continue
				}
				res[y][x] = cur
			}
		}
	}
	return res
}
