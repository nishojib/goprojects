package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// A Bookworm contains the list of books on a bookworm's shelf.
type Bookworm struct {
	Name  string `json:"name"`
	Books []Book `json:"books"`
}

// Book describes a book on a bookworm's shelf.
type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

// loadBookworms reads the file and returns the list of bookworms, and their beloved books, found therein.
func loadBookworms(filePath string) ([]Bookworm, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var bookworms []Bookworm

	buffedReader := bufio.NewReaderSize(f, 1024*1024)
	err = json.NewDecoder(buffedReader).Decode(&bookworms)
	if err != nil {
		return nil, err
	}

	return bookworms, nil
}

// findCommonBooks returns books that are on more than one bookworm's shelf.
func findCommonBooks(bookworms []Bookworm) []Book {
	booksOnShelves := booksCount(bookworms)

	var commonBooks []Book

	for book, count := range booksOnShelves {
		if count > 1 {
			commonBooks = append(commonBooks, book)
		}
	}

	return sortBooks(commonBooks)
}

// booksCount registers all the books and their occurrences from the bookworms shelves.
func booksCount(bookworms []Bookworm) map[Book]uint {
	count := make(map[Book]uint)

	for _, bookworm := range bookworms {
		for _, book := range bookworm.Books {
			count[book]++
		}
	}

	return count
}

// displayBooks prints out the titles and authors of a list of books
func displayBooks(books []Book) {
	for _, book := range books {
		fmt.Println("-", book.Title, "by", book.Author)
	}
}

// Books is a list of Books. Defining a custom type to implement sort.Interface
type byAuthor []Book

// Len implements sort.Interface by returning the length of the collection.
func (b byAuthor) Len() int { return len(b) }

// Swap implements sort.Interface and swaps two books
func (b byAuthor) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Less implements sort.Interface and returns books sorted by Author and then Title.
func (b byAuthor) Less(i, j int) bool {
	if b[i].Author != b[j].Author {
		return b[i].Author < b[j].Author
	}

	return b[i].Title < b[j].Title
}

// sortBooks sorts the books by Author and then Title.
func sortBooks(books []Book) []Book {
	sort.Sort(byAuthor(books))
	return books
}
