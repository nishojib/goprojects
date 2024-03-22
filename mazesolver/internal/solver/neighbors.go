package solver

import "image"

// neighbors returns an array of the 4 neighbors of a pixel.
// Some returned positions may be outside the maze.
func neighbors(p image.Point) []image.Point {
	return []image.Point{
		{p.X, p.Y + 1},
		{p.X, p.Y - 1},
		{p.X + 1, p.Y},
		{p.X - 1, p.Y},
	}
}
