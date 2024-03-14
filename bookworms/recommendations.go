package main

import (
	"fmt"
	"math"
	"sort"
)

// creates a set of books
type set map[Book]struct{}

// Contains implements Contains method on the set
func (s set) Contains(b Book) bool {
	_, ok := s[b]
	return ok
}

// newSet converts a list of books to set
func newSet(books []Book) set {
	s := make(set)
	for _, book := range books {
		s[book] = struct{}{}
	}
	return s
}

// Recommendation describes a recommendation for a particular book with score
type Recommendation struct {
	Book  Book
	Score float64
}

// recommend takes the list of all bookworms and a target and returns `n` number of recommendations
func recommend(allReaders []Bookworm, target Bookworm, n int) []Recommendation {
	read := newSet(target.Books)

	recommendations := map[Book]float64{}
	for _, reader := range allReaders {
		if reader.Name == target.Name {
			continue
		}

		var similarity float64
		for _, book := range reader.Books {
			if read.Contains(book) {
				similarity++
			}
		}

		if similarity == 0 {
			continue
		}

		score := math.Log(similarity) + 1
		for _, book := range reader.Books {
			if !read.Contains(book) {
				recommendations[book] += score
			}
		}
	}

	recs := bookRecommendationToListOfBooks(recommendations)

	if n > len(recs) {
		return recs
	}

	return recs[:n]
}

// bookRecommendationToListOfBooks converts a map of recommendations to a list sorted by score in desc order
func bookRecommendationToListOfBooks(rec map[Book]float64) []Recommendation {
	recs := make([]Recommendation, 0, len(rec))
	for book, score := range rec {
		recs = append(recs, Recommendation{Book: book, Score: score})
	}

	return sortRecommendations(recs)
}

// displayRecommendations prints out the titles and authors of a list of recommendations with their score
func displayRecommendations(recs []Recommendation) {
	for _, rec := range recs {
		fmt.Println("-", rec.Book.Title, "by", rec.Book.Author, "with score", rec.Score)
	}
}

// Recommendations is a list of recommendations. Defining a custom type to implement sort.Interface
type byScore []Recommendation

// Len implements sort.Interface by returning the length of the collection
func (b byScore) Len() int { return len(b) }

// Swap implements sort.Interface and swaps two recommendations
func (b byScore) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Less implements sort.Interface and returns recommendations sorted by Score
func (b byScore) Less(i, j int) bool {
	return b[i].Score > b[j].Score
}

// sortRecommendations sorts the recommendations by Score
func sortRecommendations(recs []Recommendation) []Recommendation {
	sort.Sort(byScore(recs))
	return recs
}
