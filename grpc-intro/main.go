package main

import "fmt"

type point struct {
	x, y int
}

func (p point) Move(dx, dy int) {
	p.x += dx
	p.y += dy
}

func main() {
	fmt.Println("Hello from Week 1 gRPC Chat")
	test()
}

func test() {

	var x int = 5
	var name string = "Alice"

	age := 30

	var a, b int = 10, 20

	fmt.Println(x, name, age, a, b)

	fmt.Println(swap(3, 4))

	pointuse()
}

func pointuse() {
	p := point{x: 4, y: 4}
	fmt.Println(p)
}

func swap(a int, b int) (int, int) {
	return b, a
}
