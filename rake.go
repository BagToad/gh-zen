package main

// Direction the rake is facing.
type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

// Rake is the player-controlled tool.
type Rake struct {
	X, Y    int
	Dir     Direction
	Flipped bool // false = spokes forward, true = flat forward
}

func newRake(x, y int) Rake {
	return Rake{X: x, Y: y, Dir: DirRight}
}

// cells returns the 3 positions the rake occupies (perpendicular to direction).
func (r Rake) cells() [3][2]int {
	return perpCells(r.X, r.Y, r.Dir)
}

// perpCells returns 3 positions spanning perpendicular to the given direction.
func perpCells(cx, cy int, dir Direction) [3][2]int {
	switch dir {
	case DirUp, DirDown:
		return [3][2]int{{cx - 1, cy}, {cx, cy}, {cx + 1, cy}}
	default:
		return [3][2]int{{cx, cy - 1}, {cx, cy}, {cx, cy + 1}}
	}
}

// occupies returns true if the rake covers the given cell.
func (r Rake) occupies(x, y int) bool {
	for _, c := range r.cells() {
		if c[0] == x && c[1] == y {
			return true
		}
	}
	return false
}

// spokeRune returns the character for the spoked end.
func spokeRune(d Direction) rune {
	switch d {
	case DirUp:
		return '╥'
	case DirDown:
		return '╨'
	case DirLeft:
		return '╡'
	case DirRight:
		return '╞'
	}
	return '+'
}

// flatRune returns the character for the flat end.
func flatRune(d Direction) rune {
	switch d {
	case DirUp, DirDown:
		return '━'
	default:
		return '┃'
	}
}

// rakeRune returns the display character based on direction and flip state.
func rakeRune(d Direction, flipped bool) rune {
	if flipped {
		return flatRune(d)
	}
	return spokeRune(d)
}

// delta returns the (dx, dy) for a direction.
func (d Direction) delta() (int, int) {
	switch d {
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	}
	return 0, 0
}

// rakedCell returns the raked cell type for a given direction.
func rakedCell(d Direction) Cell {
	switch d {
	case DirLeft, DirRight:
		return CellRakedH
	default:
		return CellRakedV
	}
}
