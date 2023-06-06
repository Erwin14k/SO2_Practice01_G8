package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"log"
	"fmt"
)

type Process struct {
	Pid     int    `json:"pid"`
	Nombre  string `json:"nombre"`
	Usuario string `json:"usuario"`
	Estado  string `json:"estado"`
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

type RAMInfo struct {
	TotalRAM    int  `json:"totalram"`
	RAMLibre    int  `json:"ramlibre"`
	RAMOcupada  int  `json:"ramocupada"`
}

type general struct {
	TotalRAM    int  `json:"totalram"`
	RAMLibre    int  `json:"ramlibre"`
	RAMOcupada  int  `json:"ramocupada"`
	TotalCPU    int  `json:"totalcpu"`
}

type counters struct {
	Running  int  `json:"running"`
	Sleeping int  `json:"sleeping"`
	Stopped  int  `json:"stopped"`
	Zombie   int  `json:"zombie"`
	Total    int  `json:"total"`
}

type AllData struct {
	AllGenerales    []general   `json:"AllGenerales"`
	AllTipoProcesos []Process   `json:"AllTipoProcesos"`
	AllProcesos     []counters	`json:"AllProcesos"`
}

// Function to create data by reading CPU and RAM information
func createData() (string, error) {

	// Read RAM information from "/proc/mem_grupo8"
	outRAM, err := ioutil.ReadFile("/proc/mem_grupo8")
	if err != nil {
		fmt.Println(err)
	}

	// Read CPU information from "/proc/cpu_grupo8"
	outCPU, err := ioutil.ReadFile("/proc/cpu_grupo8")
	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal CPU information into CPUInfo struct
	var cpuInfo CPUInfo
	err = json.Unmarshal(outCPU, &cpuInfo)
	if err != nil {
		fmt.Println("Error: Cpu json unmarshal failed", err)
		return "", err
	}

	// Iterate over CPU tasks and retrieve username for each UID
	for i, task := range cpuInfo.Tasks {
		uid, err := strconv.Atoi(task.Usuario)
		if err != nil {
			fmt.Println("Error: Failed to convert UID to int", err)
			return "", err
		}

		// Execute shell command to retrieve username for UID
		cmdUsr := exec.Command("sh", "-c", "grep -m 1 '"+strconv.Itoa(uid)+":' /etc/passwd | cut -d: -f1")
		outUsr, err := cmdUsr.Output()
		if err != nil {
			fmt.Println("Error: Failed to get username for UID ", task.Usuario, err)
			return "", err
		}
		username := strings.TrimSpace(string(outUsr))
		cpuInfo.Tasks[i].Usuario = username
	}

	// Unmarshal RAM information into RAMInfo struct
	var ramInfo RAMInfo
	err = json.Unmarshal(outRAM, &ramInfo)
	if err != nil {
		fmt.Println("Error: Ram json unmarshal failed", err)
		return "", err
	}

	// Create AllData struct with all the information
	allData := AllData{
		AllGenerales: []general{
			{
				TotalRAM:     ramInfo.TotalRAM,
				RAMLibre:     ramInfo.RAMLibre,
				RAMOcupada:   ramInfo.RAMOcupada,
				TotalCPU:     cpuInfo.TotalCPU,
			},
		},
		AllTipoProcesos: cpuInfo.Tasks,
		AllProcesos: []counters{
			{
				Running:  cpuInfo.Running,
				Sleeping: cpuInfo.Sleeping,
				Stopped:  cpuInfo.Stopped,
				Zombie:   cpuInfo.Zombie,
				Total:    cpuInfo.Total,
			},
		},
	}

	// Marshal AllData struct into JSON format
	allDataJSON, err := json.Marshal(allData)
	if err != nil {
		fmt.Println("Error: AllData json marshal failed", err)
		return "", err
	}

	return string(allDataJSON), nil
}

// Handler for GET requests
func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet { 
		w.WriteHeader(http.StatusMethodNotAllowed) 
		return
	}

	allData, err := createData() 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) 
		return
	}

	w.Header().Set("Content-Type", "application/json") 
	fmt.Fprint(w, allData)                              
}

// Handler for POST requests to delete a process
func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { 
		w.WriteHeader(http.StatusMethodNotAllowed) 
		return
	}

	body, err := ioutil.ReadAll(r.Body) 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) 
		return
	}

	pid, err := strconv.Atoi(string(body)) 
	if err != nil {
		w.WriteHeader(http.StatusBadRequest) 
		fmt.Fprintln(w, "Error: Invalid PID")
		return
	}

	cmd := exec.Command("sudo", "kill", strconv.Itoa(pid)) 
	err = cmd.Run()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError) 
		fmt.Fprintln(w, "Error: It is not possible to kill the process")
		return
	}

	fmt.Println("Information: Process with PID", pid, "has been deleted") 
	w.WriteHeader(http.StatusOK)              
	fmt.Fprintln(w, "Process deleted")        
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
