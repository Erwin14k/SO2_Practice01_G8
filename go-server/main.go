package main

import (
	"encoding/json" // Package for JSON encoding and decoding
	"io/ioutil"     // Package for file I/O operations
	"net/http"      // Package for HTTP-related operations
	"log"          // Package for logging
	"fmt"          // Package for formatted I/O
)

var singleModule = "/proc/mem_grupo8" // File path for the data source
var D data // Variable to hold the parsed JSON data

type process struct {
    pid        int    `json:"pid"` // Process ID
    name       string `json:"name"` // Process name
    user       string `json:"user"` // User associated with the process
    status     int    `json:"status"` // Process status
    ram        int    `json:"ram"` // RAM usage of the process
    children   []struct {
        Pid    int    `json:"Pid"` // Child process ID
        Nombre string `json:"Nombre"` // Child process name
    } `json:"children"` // Child processes of the main process
}

type data struct {
    processes []process   `json:"processes"` // List of processes
    total_ram     string  `json:"total_ram"` // Total RAM available
    free_ram      string  `json:"free_ram"` // Free RAM
    ram_occupied  string  `json:"ram_occupied"` // Occupied RAM
    counters struct {
        running   int `json:"running"` // Number of running processes
        suspended int `json:"suspended"` // Number of suspended processes
        stopped   int `json:"stopped"` // Number of stopped processes
        zombies   int `json:"zombies"` // Number of zombie processes
        total     int `json:"total"` // Total number of processes
    } `json:"counters"` // Process counters
}

func getProcessData() {
    // Read data from the specified file
    data, err := ioutil.ReadFile(singleModule)
    if err != nil {
        fmt.Println(err)
    }
    // Unmarshal the JSON data into the 'D' variable
    err = json.Unmarshal(data, &D)
    if err != nil {
        fmt.Println(err)
    }
}

func createData() string {
	getProcessData() // Call the function to retrieve process data

	b, err := json.Marshal(D) // Convert the data into JSON format
	if err != nil {
		log.Println("Error converting to JSON") // Log an error if the conversion fails
		return ""
	}
	return string(b) // Return the JSON data as a string
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { // Check if the HTTP method is GET
		w.WriteHeader(http.StatusMethodNotAllowed) // Return HTTP 405 Method Not Allowed if it's not GET
		return
	}

	allData := createData() // Retrieve all the process data as JSON
	if allData == "" {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if the data is empty
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the response header to indicate JSON content type
	fmt.Fprint(w, allData) // Write the JSON data to the response writer
}


func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // Check if the HTTP method is POST
		w.WriteHeader(http.StatusMethodNotAllowed) // Return HTTP 405 Method Not Allowed if it's not POST
		return
	}

	body, err := ioutil.ReadAll(r.Body) // Read the request body
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Return HTTP 500 Internal Server Error if there's an error reading the body
		return
	}

	fmt.Println("Information: Process with PID", string(body), "has been deleted") // Print the information about the deleted process
	w.WriteHeader(http.StatusOK) // Set HTTP 200 OK status code
	fmt.Fprintln(w, "Process deleted") // Write the response message to the response writer
}

func main() {
	fmt.Println("************************************************************")
	fmt.Println("*                 SO2 Practica 1 - Grupo 8                 *")
	fmt.Println("************************************************************")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(w, r) // Call handleGet for GET requests
		case http.MethodPost:
			handlePost(w, r) // Call handlePost for POST requests
		default:
			w.WriteHeader(http.StatusMethodNotAllowed) // Return HTTP 405 Method Not Allowed for other methods
		}
	})

	fmt.Println("Server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil)) // Start the server on port 8080
}
