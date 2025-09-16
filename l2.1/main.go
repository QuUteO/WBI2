package main

import "fmt"

// Что выведет этот код?
func main() {
	a := [5]int{76, 77, 78, 79, 80}
	// создается новый слайс, который будет выглядеть так [77,78,79]
	var b []int = a[1:4]
	fmt.Println(b)
}
