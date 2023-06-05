package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type Process struct {
	Pid     int    `json:"pid"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Estado  int    `json:"estado"`
	Ram     int    `json:"ram"`
	Padre   int    `json:"padre"`
}

type CPUInfo struct {
	TotalCPU int       `json:"totalcpu"`
	Running  int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped  int       `json:"stopped"`
	Zombie   int       `json:"zombie"`
	Total    int       `json:"total"`
	Tasks    []Process `json:"tasks"`
}

func createData() (string, error) {
	cmdRAM := exec.Command("sh", "-c", "cat /proc/mem_grupo8")
	outRAM, err := cmdRAM.CombinedOutput()
	if err != nil {
		fmt.Println("Error: Ram file cannot be readed", err)
		return "", err
	}

	cmdCPU := exec.Command("sh", "-c", "cat /proc/cpu_grupo8")
	outCPU, err := cmdCPU.CombinedOutput()
	if err != nil {
		fmt.Println("Error: Cpu file cannot be readed", err)
		return "", err
	}

	var cpuInfo CPUInfo
	err = json.Unmarshal([]byte(outCPU), &cpuInfo)
	if err != nil {
		fmt.Println("Error: Cpu json unmarshal failed", err)
		return "", err
	}

	for i, task := range cpuInfo.Tasks {
		uid, err := strconv.Atoi(task.Usuario)
		if err != nil {
			fmt.Println("Error: Failed to convert UID to int", err)
			return "", err
		}

		cmdUsr := exec.Command("sh", "-c", "grep -m 1 '"+strconv.Itoa(uid)+":' /etc/passwd | cut -d: -f1")
		outUsr, err := cmdUsr.Output()
		if err != nil {
			fmt.Println("Error: Failed to get username for UID ", task.Usuario, err)
			return "", err
		}
		username := strings.TrimSpace(string(outUsr))
		cpuInfo.Tasks[i].Usuario = username
	}

	cpuData, err := json.Marshal(cpuInfo)
	if err != nil {
		fmt.Println("Error: Cpu json marshal failed", err)
		return "", err
	}

	// --------- RAM ---------
	var mapRAM map[string]int
	err = json.Unmarshal([]byte(outRAM), &mapRAM)
	if err != nil {
		fmt.Println("Error: Ram json unmarshal failed", err)
		return "", err
	}

	ramData, err := json.Marshal(mapRAM)
	if err != nil {
		fmt.Println("Error: Ram json marshal failed", err)
		return "", err
	}

	allData := fmt.Sprintf(`{"cpuData": %s, "ramData": %s}`, cpuData, ramData)
	return allData, nil
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { // Comprobar si el método HTTP es GET
		w.WriteHeader(http.StatusMethodNotAllowed) // Devolver HTTP 405 Method Not Allowed si no es GET
		return
	}

	allData, err := createData() // Obtener todos los datos de los procesos en formato JSON
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) // Devolver HTTP 500 Internal Server Error si los datos están vacíos
		return
	}

	w.Header().Set("Content-Type", "application/json") // Establecer la cabecera de la respuesta para indicar el tipo de contenido JSON
	fmt.Fprint(w, allData)                              // Escribir los datos JSON en el escritor de la respuesta
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
