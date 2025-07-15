{
    "code": "package main

import (
    \"encoding/json\"
    \"fmt\"
    \"io/ioutil\"
    \"log\"
    \"net/http\"

    \"github.com/go-playground/validator/v10\"
    \"github.com/gorilla/mux\"
)

type Book struct {
    ID     string `json:\"id\"`
    Title  string `json:\"title\" validate:\"required\"`
    Author string `json:\"author\" validate:\"required\"`
    Year   int    `json:\"year\" validate:\"required,min=1000,max=2100\"`
}

var books []Book
var validate *validator.Validate

func main() {
    validate = validator.New()
    r := mux.NewRouter()

    r.HandleFunc(\"/books\", getBooks).Methods(\"GET\")
    r.HandleFunc(\"/books\", createBook).Methods(\"POST\")
    r.HandleFunc(\"/books/{id}\", getBook).Methods(\"GET\")
    r.HandleFunc(\"/books/{id}\", updateBook).Methods(\"PUT\")
    r.HandleFunc(\"/books/{id}\", deleteBook).Methods(\"DELETE\")

    fmt.Println(\"Server started at http://localhost:8000\")
    log.Fatal(http.ListenAndServe(\":8000\", r))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set(\"Content-Type\", \"application/json\")
    json.NewEncoder(w).Encode(books)
}

func createBook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set(\"Content-Type\", \"application/json\")

    var book Book
    body, _ := ioutil.ReadAll(r.Body)
    err := json.Unmarshal(body, &book)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, \"Error: %v\", err)
        return
    }

    err = validate.Struct(book)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        fmt.Fprintf(w, \"Error: %v\", err)
        return
    }

    book.ID = fmt.Sprintf(\"%d\", len(books)+1)
    books = append(books, book)
    json.NewEncoder(w).Encode(book)
}

func getBook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set(\"Content-Type\", \"application/json\")
    params := mux.Vars(r)
    id := params[\"id\"]

    for _, book := range books {
        if book.ID == id {
            json.NewEncoder(w).Encode(book)
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, \"Book not found\")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set(\"Content-Type\", \"application/json\")
    params := mux.Vars(r)
    id := params[\"id\"]

    for i, book := range books {
        if book.ID == id {
            var updatedBook Book
            body, _ := ioutil.ReadAll(r.Body)
            err := json.Unmarshal(body, &updatedBook)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprintf(w, \"Error: %v\", err)
                return
            }

            err = validate.Struct(updatedBook)
            if err != nil {
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprintf(w, \"Error: %v\", err)
                return
            }

            updatedBook.ID = id
            books[i] = updatedBook
            json.NewEncoder(w).Encode(updatedBook)
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, \"Book not found\")
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
    w.Header().Set(\"Content-Type\", \"application/json\")
    params := mux.Vars(r)
    id := params[\"id\"]

    for i, book := range books {
        if book.ID == id {
            books = append(books[:i], books[i+1:]...)
            json.NewEncoder(w).Encode(books)
            return
        }
    }

    w.WriteHeader(http.StatusNotFound)
    fmt.Fprintf(w, \"Book not found\")
}",

    "tests": "package main

import (
    \"bytes\"
    \"encoding/json\"
    \"net/http\"
    \"net/http/httptest\"
    \"testing\"

    \"github.com/gorilla/mux\"
)

func TestGetBooks(t *testing.T) {
    req, err := http.NewRequest(\"GET\", \"/books\", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(getBooks)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf(\"Handler returned wrong status code: got %v want %v\", status, http.StatusOK)
    }
}

func TestCreateBook(t *testing.T) {
    book := Book{Title: \"Test Book\", Author: \"Test Author\", Year: 2022}
    body, _ := json.Marshal(book)

    req, err := http.NewRequest(\"POST\", \"/books\", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    handler := http.HandlerFunc(createBook)
    handler.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf(\"Handler returned wrong status code: got %v want %v\", status, http.StatusOK)
    }

    var createdBook Book
    err = json.Unmarshal(rr.Body.Bytes(), &createdBook)
    if err != nil {
        t.Fatal(err)
    }

    if createdBook.Title != book.Title || createdBook.Author != book.Author || createdBook.Year != book.Year {
        t.Errorf(\"Book data mismatch: got %+v want %+v\", createdBook, book)
    }
}

func TestGetBook(t *testing.T) {
    req, err := http.NewRequest(\"GET\", \"/books/1\", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    r := mux.NewRouter()
    r.HandleFunc(\"/books/{id}\", getBook)
    r.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf(\"Handler returned wrong status code: got %v want %v\", status, http.StatusOK)
    }
}

func TestUpdateBook(t *testing.T) {
    book := Book{Title: \"Test Book\", Author: \"Test Author\", Year: 2022}
    body, _ := json.Marshal(book)

    req, err := http.NewRequest(\"PUT\", \"/books/1\", bytes.NewBuffer(body))
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    r := mux.NewRouter()
    r.HandleFunc(\"/books/{id}\", updateBook)
    r.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf(\"Handler returned wrong status code: got %v want %v\", status, http.StatusOK)
    }

    var updatedBook Book
    err = json.Unmarshal(rr.Body.Bytes(), &updatedBook)
    if err != nil {
        t.Fatal(err)
    }

    if updatedBook.Title != book.Title || updatedBook.Author != book.Author || updatedBook.Year != book.Year {
        t.Errorf(\"Book data mismatch: got %+v want %+v\", updatedBook, book)
    }
}

func TestDeleteBook(t *testing.T) {
    req, err := http.NewRequest(\"DELETE\", \"/books/1\", nil)
    if err != nil {
        t.Fatal(err)
    }

    rr := httptest.NewRecorder()
    r := mux.NewRouter()
    r.HandleFunc(\"/books/{id}\", deleteBook)
    r.ServeHTTP(rr, req)

    if status := rr.Code; status != http.StatusOK {
        t.Errorf(\"Handler returned wrong status code: got %v want %v\", status, http.StatusOK)
    }
}",

    "documentation": "This is a simple REST API implementation in Go for managing a list of books. It provides CRUD operations (Create, Read, Update, Delete) for books. The API uses the gorilla/mux package for routing and go-playground/validator for input validation. The Book struct represents a book with fields like ID, Title, Author, and Year. The API handles requests for getting all books, creating a new book, getting a specific book, updating a book, and deleting a book. The code includes proper error handling and input validation using the validator package. The tests directory contains unit tests for the API endpoints, covering different scenarios like getting books, creating a book, getting a specific book, updating a book, and deleting a book.",

    "dependencies": [
        "github.com/go-playground/validator/v10",
        "github.com/gorilla/mux"
    ]
}