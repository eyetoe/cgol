package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"
)

const width int = 80
const height int = 80
const prob int = 60

// MAIN ///////////////////////////////////////////
func main() {

	new := newWorld(prob)
	new.display()
	cont()
	for i := 0; true; i++ {
		new = new.step()
		new.display()
		fmt.Println("iteration:", i)
		time.Sleep(200 * time.Millisecond)
	}

	//fmt.Println(checkAlive(false, 3))

	//		time.Sleep(1000 * time.Millisecond)
}

func myNeighbors(c coords) []coords {
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

func (w world) step() world {

	// zero out neighbor values
	for x, _ := range w {
		for y, _ := range w[x] {
			w[x][y].neighbors = 0
		}
	}
	// For each cell assign neighbor cells +1
	for xi, x := range w {
		for yi, _ := range x {
			if w[xi][yi].alive == true {
				nbrs := myNeighbors(coords{xi, yi})
				//fmt.Println(coords{xi, yi}, "neighbors are:", nbrs)

				for _, n := range nbrs {
					if n.x >= 0 && n.x < width && n.y >= 0 && n.y < height {
						w[n.x][n.y].neighbors = w[n.x][n.y].neighbors + 1

					}
				}
			}
		}
	}
	// set dead or alive
	for xi, x := range w {
		for yi, _ := range x {
			// set alive here
			//fmt.Println(w[xi][yi])
			w[xi][yi].alive = checkAlive(w[xi][yi].alive, w[xi][yi].neighbors)
			//fmt.Printf("starting vals: %t, %d => new val: %t\n", w[xi][yi].alive, w[xi][yi].neighbors, checkAlive(w[xi][yi].alive, w[xi][yi].neighbors))
		}
	}
	return w
}

func checkAlive(alive bool, neighbors int) bool {
	switch alive {
	case true:
		switch {
		case neighbors < 2:
			return false
		case neighbors > 1 && neighbors < 4:
			return true
		case neighbors > 3:
			return false
		}
	case false:
		switch {
		case neighbors == 3:
			return true
		default:
			return false
		}
	default:
		return false
	}
	return false
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

func (w world) display() {
	clearscreen()
	// debug
	//	for xi, x := range w {
	//		for yi, y := range x {
	//			fmt.Printf("%d,%d %t	", xi, yi, y.alive)
	//		}
	//		fmt.Println()
	//	}
	// end debug
	for _, x := range w {
		for _, y := range x {
			switch y.alive {
			case true:
				fmt.Printf("[]")
			case false:
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
}

type coords struct {
	x int
	y int
}

type cell struct {
	alive     bool
	neighbors int
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
