package solver

import "image/color"

// palette contains the colors of the different types of pixels in our maze.
type palette struct {
	wall     color.RGBA
	path     color.RGBA
	entrance color.RGBA
	treasure color.RGBA
	solution color.RGBA
	explored color.RGBA
}

// defaultPalette returns the color palette of our maze.
func defaultPalette() palette {
	return palette{
		wall:     color.RGBA{R: 0, G: 0, B: 0, A: 255},
		path:     color.RGBA{R: 255, G: 255, B: 255, A: 255},
		entrance: color.RGBA{R: 0, G: 191, B: 255, A: 255},
		treasure: color.RGBA{R: 255, G: 0, B: 128, A: 255},
		solution: color.RGBA{R: 225, G: 140, B: 0, A: 255},
		explored: color.RGBA{R: 0, G: 128, B: 255, A: 255},
	}
}
