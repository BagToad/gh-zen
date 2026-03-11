package main

import "math/rand"

// Cell represents a single tile in the garden grid.
type Cell int

const (
	CellEmpty     Cell = iota // outside playable area
	CellSand                  // unraked sand ░
	CellRakedH                // horizontally raked ─
	CellRakedV                // vertically raked │
	CellFlattened             // fully flattened ·
	CellRock                  // part of a rock obstacle
)

// Garden holds the grid state.
type Garden struct {
	Width   int
	Height  int
	Cells   [][]Cell
	Texture [][]int // random per-cell texture seed
}

func newGarden(w, h int) *Garden {
	g := &Garden{Width: w, Height: h}
	g.Cells = make([][]Cell, h)
	g.Texture = make([][]int, h)
	for y := range g.Cells {
		g.Cells[y] = make([]Cell, w)
		g.Texture[y] = make([]int, w)
		for x := range g.Cells[y] {
			g.Cells[y][x] = CellSand
			g.Texture[y][x] = rand.Intn(100)
		}
	}
	return g
}

func (g *Garden) inBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

func (g *Garden) at(x, y int) Cell {
	if !g.inBounds(x, y) {
		return CellEmpty
	}
	return g.Cells[y][x]
}

func (g *Garden) set(x, y int, c Cell) {
	if g.inBounds(x, y) {
		g.Cells[y][x] = c
	}
}

func (g *Garden) isPassable(x, y int) bool {
	if !g.inBounds(x, y) {
		return false
	}
	c := g.Cells[y][x]
	return c != CellRock && c != CellEmpty
}

func (g *Garden) isSand(x, y int) bool {
	if !g.inBounds(x, y) {
		return false
	}
	return g.Cells[y][x] == CellSand
}

// countUnraked returns how many sand cells remain unraked.
func (g *Garden) countUnraked() int {
	count := 0
	for y := range g.Cells {
		for _, c := range g.Cells[y] {
			if c == CellSand {
				count++
			}
		}
	}
	return count
}

// placeRocks randomly places multi-char rocks, keeping a clear zone around the rake spawn.
func (g *Garden) placeRocks(rakeX, rakeY int) {
	shapes := rockShapes()
	// fewer but larger rocks
	area := g.Width * g.Height
	numRocks := area / 400
	if numRocks < 1 {
		numRocks = 1
	}
	if numRocks > 4 {
		numRocks = 4
	}

	clearRadius := 7
	attempts := 0
	placed := 0
	for placed < numRocks && attempts < numRocks*50 {
		attempts++
		shape := shapes[rand.Intn(len(shapes))]
		// random position with margin from edges
		x := rand.Intn(g.Width-shape.w-2) + 1
		y := rand.Intn(g.Height-shape.h-2) + 1

		if g.canPlaceRock(x, y, shape, rakeX, rakeY, clearRadius) {
			g.stampRock(x, y, shape)
			placed++
		}
	}
}

func (g *Garden) canPlaceRock(x, y int, s rockShape, rakeX, rakeY, clearR int) bool {
	for dy := 0; dy < s.h; dy++ {
		for dx := 0; dx < s.w; dx++ {
			if s.mask[dy][dx] == 0 {
				continue
			}
			px, py := x+dx, y+dy
			if !g.inBounds(px, py) {
				return false
			}
			if g.Cells[py][px] != CellSand {
				return false
			}
			// check clear zone around rake spawn
			if abs(px-rakeX) <= clearR && abs(py-rakeY) <= clearR {
				return false
			}
		}
	}
	return true
}

func (g *Garden) stampRock(x, y int, s rockShape) {
	for dy := 0; dy < s.h; dy++ {
		for dx := 0; dx < s.w; dx++ {
			if s.mask[dy][dx] == 1 {
				g.Cells[y+dy][x+dx] = CellRock
			}
		}
	}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
