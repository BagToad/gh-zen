package main

// rockShape defines a multi-character rock via a bitmask.
type rockShape struct {
	w, h int
	mask [][]int
}

func rockShapes() []rockShape {
	return []rockShape{
		// Massive boulder 8×5
		{w: 8, h: 5, mask: [][]int{
			{0, 0, 1, 1, 1, 1, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 0},
			{1, 1, 1, 1, 1, 1, 1, 1},
			{0, 1, 1, 1, 1, 1, 1, 0},
			{0, 0, 1, 1, 1, 1, 0, 0},
		}},
		// Wide slab 10×3
		{w: 10, h: 3, mask: [][]int{
			{0, 1, 1, 1, 1, 1, 1, 1, 1, 0},
			{1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			{0, 0, 1, 1, 1, 1, 1, 1, 0, 0},
		}},
		// Tall pillar 4×7
		{w: 4, h: 7, mask: [][]int{
			{0, 1, 1, 0},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{1, 1, 1, 1},
			{0, 1, 1, 0},
		}},
		// Irregular crag 7×5
		{w: 7, h: 5, mask: [][]int{
			{0, 0, 1, 1, 1, 0, 0},
			{0, 1, 1, 1, 1, 1, 0},
			{1, 1, 1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1, 0, 0},
			{0, 1, 1, 0, 0, 0, 0},
		}},
		// Chunky 5×4
		{w: 5, h: 4, mask: [][]int{
			{0, 1, 1, 1, 0},
			{1, 1, 1, 1, 1},
			{1, 1, 1, 1, 1},
			{0, 0, 1, 1, 0},
		}},
	}
}

// rockChar returns the display character for a rock cell based on its
// neighbours (up, down, left, right) also being rock cells.
func rockChar(up, down, left, right bool) rune {
	switch {
	case !up && !down && !left && !right:
		return '●'
	case up && down && left && right:
		return '█'
	case !up && down && !left && right:
		return '▛'
	case !up && down && left && !right:
		return '▜'
	case up && !down && !left && right:
		return '▙'
	case up && !down && left && !right:
		return '▟'
	case up && down:
		return '█'
	case left && right:
		return '█'
	default:
		return '█'
	}
}
