package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

type Book struct {
    ID     uint   `json:"id" gorm:"primary_key"`
    Title  string `json:"title"`
    Author string `json:"author"`
}

var db *gorm.DB
var err error

func setupDB() {
    db, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
    if err != nil {
        log.Fatal(err)
    }

    db.AutoMigrate(&Book{})
}

func createBook(c *gin.Context) {
    var book Book
    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.Create(&book)
    c.JSON(http.StatusCreated, book)
}

func getBooks(c *gin.Context) {
    var books []Book
    db.Find(&books)
    c.JSON(http.StatusOK, books)
}

func getBook(c *gin.Context) {
    var book Book
    id := c.Param("id")
    db.First(&book, id)
    if book.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }
    c.JSON(http.StatusOK, book)
}

func updateBook(c *gin.Context) {
    var book Book
    id := c.Param("id")
    db.First(&book, id)
    if book.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    if err := c.ShouldBindJSON(&book); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    db.Save(&book)
    c.JSON(http.StatusOK, book)
}

func deleteBook(c *gin.Context) {
    var book Book
    id := c.Param("id")
    db.First(&book, id)
    if book.ID == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
        return
    }

    db.Delete(&book)
    c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

func setupRouter() *gin.Engine {
    r := gin.Default()

    r.Use(gin.Logger())
    r.Use(gin.Recovery())

    limiter := gin.NewRateLimiter(1, 1)
    r.Use(limiter)

    r.POST("/books", createBook)
    r.GET("/books", getBooks)
    r.GET("/books/:id", getBook)
    r.PUT("/books/:id", updateBook)
    r.DELETE("/books/:id", deleteBook)

    return r
}

func main() {
    setupDB()
    r := setupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    err := r.Run(fmt.Sprintf(":%s", port))
    if err != nil {
        log.Fatal(err)
    }
}