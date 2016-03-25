package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func newWorld() world {
	var new world
	for c := 1; c <= height; c++ {
		var row []cell
		for r := 1; r <= width; r++ {
			if roll(1, 100) <= prob {
				row = append(row, cell{alive: true})
			} else {
				row = append(row, cell{alive: false})
			}
		}
		new = append(new, row)
	}
	return new
}

func (w world) display() {
	for _, r := range w {
		for _, c := range r {
			switch c.alive {
			case true:
				fmt.Printf("[]")
			case false:
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
}

// MAIN ///////////////////////////////////////////
func main() {
	fmt.Println(width, height)
	smorg := newWorld()
	smorg.display()
}

type coords struct {
	x int
	y int
}

type cell struct {
	alive     bool
	neighbors int
}

func (c coords) neighbors() []coords {
	var xn = []int{c.x - 1, c.x, c.x + 1}
	var yn = []int{c.y - 1, c.y, c.y + 1}

	var neighbors []coords

	for _, a := range xn {
		for _, b := range yn {
			// reflect.DeepEqual is often wrong to use
			// but in this context there is very low variability
			// though it may still be wrong.
			if !reflect.DeepEqual(c, coords{a, b}) {
				neighbors = append(neighbors, coords{a, b})
			}
		}
	}
	return neighbors
}

const width int = 20
const height int = 20
const prob int = 10

type world [][]cell

func roll(n int, d int) int {
	num := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		num = num + r.Intn(d)
	}
	num = num / n
	return num + 1
}
