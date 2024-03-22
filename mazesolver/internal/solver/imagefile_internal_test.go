package solver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenMaze_errors(t *testing.T) {
	testCases := map[string]struct {
		input string
		err   string
	}{
		"no such file": {
			input: "nosuchfile.png",
			err:   "no such file or directory",
		},
		"not a rgba png": {
			input: "testdata/rgb.png",
			err:   "expected RGBA image, got *image.Paletted",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			img, err := openMaze(testCase.input)

			assert.Nil(t, img)
			assert.Error(t, err)
			assert.ErrorContains(t, err, testCase.err)
		})
	}
}
