package main

import (
	"reflect"
	"slices"
	"testing"
)

var (
	handmaidsTale = Book{Author: "Margaret Atwood", Title: "The Handmaid's Tale"}
	oryxAndCrake  = Book{Author: "Margaret Atwood", Title: "Oryx and Crake"}
	theBellJar    = Book{Author: "Sylvia Plath", Title: "The Bell Jar"}
	janeEyre      = Book{Author: "Charlotte Brontë", Title: "Jane Eyre"}
	villette      = Book{Author: "Charlotte Brontë", Title: "Villette"}
	ilPrincipe    = Book{Author: "Niccolò Machiavelli", Title: "Il Principe"}
)

func TestBookworms(t *testing.T) {

	type testCase struct {
		bookwormsFile string
		want          []Bookworm
		wantErr       bool
	}

	tests := map[string]testCase{
		"file exists": {
			bookwormsFile: "testdata/bookworms.json",
			want: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			wantErr: false,
		},
		"file does not exist": {
			bookwormsFile: "testdata/does-not-exist.json",
			want:          nil,
			wantErr:       true,
		},
		"invalid JSON": {
			bookwormsFile: "testdata/invalid.json",
			want:          nil,
			wantErr:       true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := loadBookworms(tc.bookwormsFile)
			if err != nil && !tc.wantErr {
				t.Fatalf("expected an error %s, got none", err.Error())
			}

			if err == nil && tc.wantErr {
				t.Fatalf("expected no error, got one %s", err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("different result: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBooksCount(t *testing.T) {
	tt := map[string]struct {
		input []Bookworm
		want  map[Book]uint
	}{
		"nominal use case": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: map[Book]uint{handmaidsTale: 2, theBellJar: 1, oryxAndCrake: 1, janeEyre: 1},
		},
		"no bookworms": {
			input: []Bookworm{},
			want:  map[Book]uint{},
		},
		"bookworm without books": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{}},
				{Name: "Peggy", Books: []Book{}},
			},
			want: map[Book]uint{},
		},
		"bookworm with twice the same book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar, handmaidsTale}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, handmaidsTale, janeEyre}},
			},
			want: map[Book]uint{handmaidsTale: 3, theBellJar: 1, oryxAndCrake: 1, janeEyre: 1},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := booksCount(tc.input)
			if !equalBooksCount(t, got, tc.want) {
				t.Fatalf("got a different list of books: %v, expected %v", got, tc.want)
			}
		})
	}
}

func TestFindCommonBooks(t *testing.T) {
	tt := map[string]struct {
		input []Bookworm
		want  []Book
	}{
		"no common book": {
			input: []Bookworm{
				{Name: "Fadi", Books: []Book{handmaidsTale, theBellJar}},
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
			},
			want: nil,
		},
		"one common book": {
			input: []Bookworm{
				{Name: "Peggy", Books: []Book{oryxAndCrake, janeEyre}},
				{Name: "Did", Books: []Book{janeEyre}},
			},
			want: []Book{janeEyre},
		},
		"three bookworms have the same books on their shelves": {
			input: []Bookworm{
				{Name: "Peggy", Books: []Book{oryxAndCrake, ilPrincipe, janeEyre}},
				{Name: "Did", Books: []Book{janeEyre}},
				{Name: "Ali", Books: []Book{janeEyre, ilPrincipe}},
			},
			want: []Book{janeEyre, ilPrincipe},
		},
		"output is sorted by authors and then title": {
			input: []Bookworm{
				{Name: "Peggy", Books: []Book{ilPrincipe, janeEyre, villette}},
				{Name: "Did", Books: []Book{janeEyre}},
				{Name: "Ali", Books: []Book{villette, ilPrincipe}},
			},
			want: []Book{janeEyre, villette, ilPrincipe},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			got := findCommonBooks(tc.input)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("got a different list of books: %v, expected %v", got, tc.want)
			}
		})
	}
}

func equalBooksCount(t *testing.T, got, want map[Book]uint) bool {
	t.Helper()

	if len(got) != len(want) {
		return false
	}

	for book, targetCount := range want {
		count, ok := got[book]
		if !ok || targetCount != count {
			return false
		}
	}

	return true
}
