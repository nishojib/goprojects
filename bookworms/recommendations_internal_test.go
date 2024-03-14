package main

import (
	"slices"
	"testing"
)

func TestRecommend(t *testing.T) {
	tt := map[string]struct {
		input  []Bookworm
		target Bookworm
		n      int
		want   []Recommendation
	}{
		"correct recommendation": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			target: Bookworm{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
			n:      1,
			want:   []Recommendation{{Book: oryxAndCrake, Score: 1}},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := recommend(tc.input, tc.target, tc.n)

			if !slices.Equal(got, tc.want) {
				t.Fatalf("got a different list of books: %v, expected %v", got, tc.want)
			}
		})
	}
}
