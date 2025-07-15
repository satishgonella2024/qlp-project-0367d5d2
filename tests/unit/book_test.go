package main

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
    setupDB()
    r := setupRouter()

    book := Book{Title: "Test Book", Author: "Test Author"}
    body, _ := json.Marshal(book)

    req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetBooks(t *testing.T) {
    setupDB()
    r := setupRouter()

    req, _ := http.NewRequest("GET", "/books", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetBook(t *testing.T) {
    setupDB()
    r := setupRouter()

    req, _ := http.NewRequest("GET", "/books/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateBook(t *testing.T) {
    setupDB()
    r := setupRouter()

    book := Book{Title: "Updated Book", Author: "Updated Author"}
    body, _ := json.Marshal(book)

    req, _ := http.NewRequest("PUT", "/books/1", bytes.NewBuffer(body))
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteBook(t *testing.T) {
    setupDB()
    r := setupRouter()

    req, _ := http.NewRequest("DELETE", "/books/1", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}