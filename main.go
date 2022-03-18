package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"unique_index"`
	PersonID   int
}


func uploadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params["name"])
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 1. parse input, type multipart/form-data.
	r.ParseMultipartForm(10 << 20)

	// 2. retreive fail from posted form-data
	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Println("Error retreaving data from the file: \n")
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Println("Uploaded file:  \n", handler.Filename)
	fmt.Println("File size:  \n", handler.Size)
	fmt.Println("File: header \n", handler.Header)

	s := strings.Split(handler.Filename, ".")

	// 3. write temporary file on our server
	tempFile, err := ioutil.TempFile("temp-images", "upload-*."+s[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	tempFile.Write(fileBytes)

	// 4. return whether or not this operation has been successful
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func sendFile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	
}

// var (
// 	person = &Person{Name: "Jack", Email: "jack@gmail.com"}
// 	books  = []Book{
// 		{Title: "Book 1", Author: "Author 1", CallNumber: 123, PersonID: 1},
// 		{Title: "Book 2", Author: "Author 1", CallNumber: 12345, PersonID: 1},
// 	}
// )

// type Articles []Article

func getPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// var people []Person
	// db.Find(&people)

	// json.NewEncoder(w).Encode(&people)

	var names []string

	files, err := ioutil.ReadDir("./temp-images/")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		names = append(names, file.Name())
	}

	json.NewEncoder(w).Encode(&names)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)

	var person Person
	var books []Book

	db.First(&person, params["id"])
	db.Model(&person).Related(&books)

	person.Books = books

	json.NewEncoder(w).Encode(person)
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	var person Person
	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := db.Create(&person)

	err := createdPerson.Error

	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(createdPerson)
	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(params["name"]))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, "./temp-images/" + params["name"])
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HomePage Endpoint Hit")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/people", getPeople).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/people/{id}", getPerson).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/people", createPerson).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/files/{name}", uploadFile).Methods("POST", "OPTIONS")
	myRouter.HandleFunc("/files", sendFile).Methods("GET", "OPTIONS")
	myRouter.HandleFunc("/downloadFile/{name}", downloadFile).Methods("GET", "OPTIONS")
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	// Loading environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	name := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, name, password, dbPort)

	// Openning connection to database
	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}
	defer db.Close()

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// db.Create(person)

	// for idx := range books {
	// 	db.Create(&books[idx])
	// }

	handleRequests()
}
