package main

import (
	"fmt"
)

func myAbs(x float64) float64 {
	if x >= 0 {
		return x
	}
	return -x
}

// Sqrt - вычисление Квадратного корня
func Sqrt(x float64) float64 {
	var xPrev float64
	var xCur = float64(1.0)
	eps := 0.0001
	for myAbs(xCur-xPrev) > eps {
		xPrev = xCur
		xCur = (0.5 * (xCur + x/xCur))
	}
	return xCur
}

func main() {
	fmt.Println(Sqrt(2))
}
