package solver

import (
	"image"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSolver_explore(t *testing.T) {
	testCases := map[string]struct {
		inputImage string
		wantSize   int
	}{
		"cross": {
			inputImage: "testdata/explore_cross.png",
			wantSize:   2,
		},
		"dead end": {
			inputImage: "testdata/explore_deadend.png",
			wantSize:   0,
		},
		"double": {
			inputImage: "testdata/explore_double.png",
			wantSize:   1,
		},
		"treasure": {
			inputImage: "testdata/explore_treasure.png",
			wantSize:   0,
		},
		"treasure only": {
			inputImage: "testdata/explore_treasureonly.png",
			wantSize:   0,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			maze, err := openMaze(testCase.inputImage)
			require.NoError(t, err)

			s := &Solver{
				maze:           maze,
				palette:        defaultPalette(),
				pathsToExplore: make(chan *path, 3),
				quit:           make(chan struct{}),
				exploredPixels: make(chan image.Point, 3),
			}

			s.explore(&path{at: image.Point{0, 2}})
			assert.Equal(t, testCase.wantSize, len(s.pathsToExplore))
		})
	}
}
