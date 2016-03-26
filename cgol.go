package main

import (
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type world [][]cell

type coords struct {
	x int
	y int
}

type cell struct {
	alive     bool
	neighbors int
}

// default token 'changeme' gets caught in main to allow for random replacement, which can't happen outside of main.
var token = flag.String("t", "changeme", "Set the characters (2 ascii characters recommended) that represent a living cell. By default the characters will be randomly selected from a set of nice ones :)")
var width = flag.Int("w", 40, "Width of the game field")
var height = flag.Int("h", 40, "Height of the game field")
var prob = flag.Int("p", 33, "Percent chance any starting cell is 'alive'")
var genesis = flag.Int("g", 200, "Respawn world after this many itterations. Like a terminal screensaver. :)")
var sleep = flag.Duration("s", 100, "Percent chance any starting cell is 'alive'")
var debug = flag.Bool("d", false, "Debug enables table values in addition to regular display. Best with small world sizes e.. 20x20")
var rainbow []func(...interface{}) string

func main() {
	flag.Parse()
	var tokenSet bool = true

	// choose a random color palate
	chooseColor := func() {
		pick := rand.Intn(4)
		switch pick {
		case 0:
			rainbow = []func(...interface{}) string{
				blue,
				black,
				white,
			}
		case 1:
			rainbow = []func(...interface{}) string{
				green,
				blue,
				cyan,
			}
		case 2:
			rainbow = []func(...interface{}) string{
				red,
				yellow,
				magenta,
			}
		case 3:
			rainbow = []func(...interface{}) string{
				white,
				yellow,
				green,
			}
		}
	}
	chooseColor()

	// The language spec clearly states "Basically no one should use strings.Compare."
	// However in this case the recommended comparison operator '==' doesn't work,
	// and this performance hit only occurs once during when starting.
	// Details here:         https://golang.org/src/strings/compare.go
	if strings.Compare(*token, "changeme") == 0 {
		tokenSet = false
		token = &rtoken[rand.Intn(len(rtoken))]
	}

	// Setup the world
	new := newWorld(*prob)
	new.display()
	fmt.Println(len(new))
	for i := 0; true; i++ {
		new = new.step()
		new.display()
		fmt.Printf("%dx%d, %d%% seed, %s sleep : iteration:%d/%d\n", *width, *height, *prob, *sleep, i, *genesis)

		time.Sleep(*sleep * time.Millisecond)
		if i == *genesis {
			new = newWorld(*prob)
			if tokenSet == false {
				token = &rtoken[rand.Intn(len(rtoken))]
			}
			// re-choose color after each genesis
			chooseColor()
			i = 0
		}
	}

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
	for y := 0; y < len(w); y++ {
		for x := 0; x < len(w[0]); x++ {
			w[y][x].neighbors = 0
		}
	}
	// for each alive cell, update neighbors
	for y := 0; y < len(w); y++ {
		for x := 0; x < len(w[0]); x++ {

			// If cell is alive, update neighbors
			if w[y][x].alive == true {
				nbrs := myNeighbors(coords{x, y})

				// and 1 to neighbors
				for _, n := range nbrs {
					if n.x >= 0 && n.x < *width && n.y >= 0 && n.y < *height {
						w[n.y][n.x].neighbors++

					}
				}
			}
		}
	}
	// set next step's alive status
	for y := 0; y < len(w); y++ {
		for x := 0; x < len(w[0]); x++ {
			// set alive here
			w[y][x].alive = checkAlive(w[y][x].alive, w[y][x].neighbors)
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
	for c := 1; c <= *height; c++ {
		var row []cell
		for r := 1; r <= *width; r++ {
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
	clearscreen()
	if *debug == true {
		for y := 0; y < len(w); y++ {
			for x := 0; x < len(w[0]); x++ {
				fmt.Printf("%d,%d %t	", y, x, w[y][x].alive)
			}
			fmt.Println()
		}
	}
	for y := 0; y < len(w); y++ {
		for x := 0; x < len(w[0]); x++ {
			switch w[y][x].alive {
			case true:
				//fmt.Printf(blue(*token))
				fmt.Printf(rainbow[rand.Intn(len(rainbow))](*token))
			case false:
				fmt.Printf("  ")
			}
		}
		fmt.Println()
	}
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

func clearscreen() {
	fmt.Printf("[H[J")
}

//rtoken[rand.Intn(len(rtoken))]
var rtoken = []string{
	"()",
	"[]",
	"<(",
	"GO",
}

var green = color.New(color.FgGreen).SprintFunc()
var blue = color.New(color.FgBlue).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()
var yellow = color.New(color.FgYellow).SprintFunc()
var cyan = color.New(color.FgCyan).SprintFunc()
var magenta = color.New(color.FgMagenta).SprintFunc()
var white = color.New(color.FgHiWhite).SprintFunc()
var black = color.New(color.FgHiBlack).SprintFunc()
