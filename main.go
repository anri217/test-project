package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var books = []Book{
	{Id: "0", Title: "Journey to the Center of the Earth", Author: "Jules Verne", Quantity: 3},
	{Id: "1", Title: "Steve Jobs", Author: "Walter Isaacson", Quantity: 2},
	{Id: "2", Title: "The Idiot", Author: "Fydor Dostoevsky", Quantity: 1},
}

type Book struct {
	Id       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing the id parameter"})
		return
	}

	book, error := findBookById(id)

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found - wrong id"})
		return
	}

	if book.Quantity == 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "this book is not avaliable any more"})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing the id parameter"})
		return
	}

	book, error := findBookById(id)

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "book not found - wrong id"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func getBookById(c *gin.Context) {
	id := c.Param("id")
	book, error := findBookById(id)

	if error != nil {
		c.IndentedJSON(http.StatusNotFound, error.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func findBookById(id string) (*Book, error) {
	for i, book := range books {
		if book.Id == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("this book is not included")
}

func createBook(c *gin.Context) {
	var newBook Book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById)
	router.POST("/books", createBook)
	router.PATCH("/checkout", checkoutBook)
	router.PATCH("/return", returnBook)
	router.Run("localhost:8080")
}
