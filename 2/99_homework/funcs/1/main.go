package main

import "fmt"

type memoizeFunction func(int, ...int) interface{}

// TODO реализовать
var fibonacci memoizeFunction
var romanForDecimal memoizeFunction

//TODO Write memoization function

func memoize(function memoizeFunction) memoizeFunction {
	res := make(map[string]interface{})
	return func(n int, ar ...int) interface{} {
		s := fmt.Sprintf("%d", n)
		for _, i := range ar {
			s += fmt.Sprintf(",%d", i)
		}
		if val, ok := res[s]; ok {
			return val
		}
		val := function(n, ar...)
		res[s] = val
		return val
	}
}

// TODO обернуть функции fibonacci и roman в memoize
func init() {
	fibonacci = memoize(func(n int, ar ...int) interface{} {
		if n < 2 {
			return n
		} else {
			return fibonacci(n-1).(int) + fibonacci(n-2).(int)
		}
	})

	romanForDecimal = memoize(func(n int, ar ...int) interface{} {
		res := ""
		num := [...]int{n / 1000, (n % 1000) / 100, (n % 100) / 10, n % 10}
		digits := [][]string{{"M"}, {"C", "D"}, {"X", "L"}, {"I", "V"}}
		for i, dig := range num {
			if dig != 0 {
				if i == 0 {
					for j := 0; j < dig; j++ {
						res += "M"
					}
				} else {
					if dig == 9 {
						res += digits[i][0] + digits[i-1][0]
					} else if dig == 4 {
						res += digits[i][0] + digits[i][1]
					} else {
						if dig >= 5 {
							res += digits[i][1]
							dig -= 5
						}
						for j := 0; j < dig; j++ {
							res += digits[i][0]
						}
					}
				}
			}
		}
		return res
	})
}

func main() {
	fmt.Println("Fibonacci(45) =", fibonacci(45).(int))
	for _, x := range []int{900, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
		14, 15, 16, 17, 18, 19, 20, 25, 30, 40, 50, 60, 69, 70, 80,
		90, 99, 100, 200, 300, 400, 500, 600, 666, 700, 800, 900,
		1000, 1009, 1444, 1666, 1945, 1997, 1999, 2000, 2008, 2010,
		2012, 2500, 3000, 3999} {
		fmt.Printf("%4d = %s\n", x, romanForDecimal(x).(string))
	}
}
