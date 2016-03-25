package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"
)

const width int = 20
const height int = 20
const prob int = 10

// MAIN ///////////////////////////////////////////
func main() {
	//fmt.Println(width, height)
	var gameOfLife world
	gameOfLife = newWorld(prob)
	gameOfLife.display()
	cont()
	for {
		gameOfLife.step()
		time.Sleep(1000 * time.Millisecond)
		gameOfLife.display()
	}
}

func newWorld(prob int) world {
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

type world [][]cell

func (w world) step() world {
	next := newWorld(0)
	// For each cell assign neighbor cells +1
	for x, _ := range w {
		for y, _ := range w[x] {
			nbrs := coords{x, y}.myNeighbors()

			//fmt.Println("list neighbors for:	", coords{x, y})
			//fmt.Println("unfiltered neighbors:	", nbrs)
			for _, n := range nbrs {
				if n.x >= 0 && n.x < width && n.y >= 0 && n.y < height {
					//fmt.Println(n.x, n.y)
					next[n.x][n.y].neighbors++
				}
			}
			//fmt.Println(next)
			//cont()
		}
	}

	// set dead or alive
	for x, _ := range w {
		for y, _ := range w[x] {
			// set alive here
			next[x][y].alive = alive(next[x][y].alive, next[x][y].neighbors)
		}
	}

	//fmt.Println("all", next)
	return next
}

func alive(alive bool, neighbors int) bool {
	if alive == true {
		switch {
		case neighbors < 2:
			return false
		case neighbors > 1 && neighbors < 4:
			return true
		case neighbors > 3:
			return false
		}
	} else if alive == false && neighbors == 3 {
		return true
	}
	return false
}

type coords struct {
	x int
	y int
}

type cell struct {
	alive     bool
	neighbors int
}

func (c coords) myNeighbors() []coords {
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

func roll(n int, d int) int {
	num := 0
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		num = num + r.Intn(d)
	}
	num = num / n
	return num + 1
}

func (w world) display() {
	clearscreen()
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

func cont() {
	fmt.Printf("Press enter to continue: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	choice := string([]byte(input)[0])
	switch choice {
	default:
		return
	}
}
func clearscreen() {
	fmt.Printf("[H[J")
}
