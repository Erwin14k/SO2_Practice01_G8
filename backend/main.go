package main

import (
	"fmt"
	"time"
	"os/exec"
	"context"
	"strings"
	"strconv"
	"encoding/json"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)

type Process struct {
	Pid      int      `json:"pid"`
	Nombre   string   `json:"nombre"`
	Usuario  int      `json:"usuario"`
	Estado   int      `json:"estado"`
	Ram      int      `json:"ram"`
	Padre      int      `json:"padre"`
}

type CPUInfo struct {
	TotalCPU int       `json:"totalcpu"`
	Running int       `json:"running"`
	Sleeping int       `json:"sleeping"`
	Stopped int       `json:"stopped"`
	Zombie int       `json:"zombie"` 
	Total int       `json:"total"`
	Tasks    []Process `json:"tasks"`
}

func main() {
	fmt.Printf(" ===== Starting V5 ======\n")
	var db, errConnection = createConnection()

	if(errConnection != nil){
		fmt.Println(errConnection)
		return
	}

	var ctx = context.Background()
	fmt.Printf(" ===== Conectado a la base de datos:======\n")
	for {	
		fmt.Println(" =====DATOS OBTENIDOS DESDE EL MODULO===== ")
		fmt.Println("")

		cmdRAM := exec.Command("sh", "-c", "cat /proc/ram_202000119")
		outRAM, err := cmdRAM.CombinedOutput()
		if err != nil {
			fmt.Println("ERR CAT RAM",err)
		}

		cmdCPU := exec.Command("sh", "-c", "cat /proc/cpu_202000119")
		outCPU, err := cmdCPU.CombinedOutput()
		if err != nil {
			fmt.Println("ERR CAT CPU",err)
		}
		
		// --------- CPU --------- 
		fmt.Println("=====CPU======")
		var cpuInfo CPUInfo
		err = json.Unmarshal([]byte(outCPU), &cpuInfo)
		if err != nil {
			fmt.Println("ERR JSON CPU",err)
			return
		}

		// fmt.Println("TotalCPU:", cpuInfo.TotalCPU)
		// fmt.Println("Running:", cpuInfo.Running)
		// fmt.Println("Sleeping:", cpuInfo.Sleeping)
		// fmt.Println("Stopped:", cpuInfo.Stopped)
		// fmt.Println("Zombie:", cpuInfo.Zombie)
		// fmt.Println("Total:", cpuInfo.Total)
		// fmt.Println("Tasks:")
		// for _, task := range cpuInfo.Tasks {
		// 	fmt.Printf("  - PID: %d, Nombre: %s, Usuario: %d, Estado: %d, Ram: %d Padre:%d\n", task.Pid, task.Nombre, task.Usuario, task.Estado,task.Ram ,task.Padre)
		// }

		// --------- RAM --------- 
		fmt.Println("=====RAM======")
		var mapRAM map[string]int
		err = json.Unmarshal([]byte(outRAM), &mapRAM)
		if err != nil {
			fmt.Println("ERR JSON RAM", err)
			return
		}

		// fmt.Println("Total RAM:", mapRAM["totalram"])
		// fmt.Println("RAM Libre:", mapRAM["ramlibre"])
		// fmt.Println("RAM Ocupada:", mapRAM["ramocupada"])

		// --------- QUERYS ---------
		err = queryDelete(ctx,db)
		if err != nil {
			fmt.Println("ERR DEL QUERY",err)
			return
		}
		err = queryAddGenerales(ctx,db,mapRAM["totalram"],mapRAM["ramlibre"],mapRAM["ramocupada"],cpuInfo.TotalCPU)
		if err != nil {
			fmt.Println("ERR ADD GENERAL QUERY",err)
			return
		}
		err = queryAddTipoProceso(ctx,db,cpuInfo.Running,cpuInfo.Sleeping,cpuInfo.Stopped,cpuInfo.Zombie,cpuInfo.Total)
		if err != nil {
			fmt.Println("ERR ADD TIPO PROCESO QUERY",err)
			return
		}
		err = queryAddCPU(ctx,db,cpuInfo)
		if err != nil {
			fmt.Println("ERR ADD CPU QUERY",err)
			return
		}

		fmt.Println("=====Termino======")
		time.Sleep(3000 * time.Millisecond)
	}
}

func createConnection() (*sql.DB,error){
	fmt.Println("=====  Intentando conectar a DB  ======")
	//db, err := sql.Open("mysql", "root:SerchiBoi502@@tcp(localhost:3306)/sopes")
	db, err := sql.Open("mysql", "root:root@tcp(35.238.173.246:3306)/sopes")

	if err != nil {
		return nil, err
	}
	
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	
	return db, nil
}

func queryDelete(ctx context.Context,db *sql.DB) error{
	
	// queryDel := "DELETE FROM procesos;"

	// _, errDelCPU := db.ExecContext(ctx, queryDel)
	// if errDelCPU != nil {
	// 	return errDelCPU
	// }

	// queryDel = "DELETE FROM generales;"

   // _, errDelGen := db.ExecContext(ctx, queryDel)
   //  if errDelGen != nil {
   //      return errDelGen
   //  } 

	queryDel := "DELETE FROM tipoprocesos;"

   _, errDelGen := db.ExecContext(ctx, queryDel)
    if errDelGen != nil {
        return errDelGen
    } 
   fmt.Println("==== Todos los datos Eliminados ====")

	return nil
}

func queryAddGenerales(ctx context.Context,db *sql.DB,totalram int, ramlibre int, ramocupada int,totalcpu int) error{
	query := "INSERT INTO generales(totalram , ramlibre , ramocupada,totalcpu ) VALUES(?,?,?,?);"
	result,err := db.ExecContext(ctx, query, totalram, ramlibre, ramocupada,totalcpu)
	if err != nil {
		return err
	}
	_,err = result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println("==== Datos Generales Insertados ====")
	return nil
}

func queryAddTipoProceso(ctx context.Context,db *sql.DB,running int, sleeping int, stopped int,zombie int,total int) error{
	query := "INSERT INTO tipoprocesos(running , sleeping , stopped,zombie,total ) VALUES(?,?,?,?,?);"
	result,err := db.ExecContext(ctx, query, running, sleeping, stopped,zombie,total)
	if err != nil {
		return err
	}
	_,err = result.LastInsertId()
	if err != nil {
		return err
	}

	fmt.Println("==== Datos TipoProcesos Insertados ====")
	return nil
}


func queryAddCPU(ctx context.Context, db *sql.DB, cpuInfo CPUInfo) error {
	baseQuery := "INSERT INTO procesos(pid , nombre , UsuarioName ,estado ,ram ,padre ) VALUES "
	valueStrings := make([]string, 0, len(cpuInfo.Tasks))
	valueArgs := make([]interface{}, 0, len(cpuInfo.Tasks)*6)

	for _, task := range cpuInfo.Tasks {
		 valueStrings = append(valueStrings, "(?,?,?,?,?,?)")

		 cmdUsr := exec.Command("sh", "-c", "grep -m 1 '"+strconv.Itoa(task.Usuario)+":' /etc/passwd | cut -d: -f1")

		
		 outUsr, err := cmdUsr.Output()
		 if err != nil {
			  fmt.Println("Err Get UID:", err)
			  return err
		 }
		 valueArgs = append(valueArgs, task.Pid, task.Nombre, string(outUsr), task.Estado, task.Ram, task.Padre)
	}

	query := baseQuery + strings.Join(valueStrings, ",")
	
	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		 return err
	}

	// Eliminar todas las filas existentes de la tabla
	_, err = tx.Exec("DELETE FROM procesos")
	if err != nil {
		 // Deshacer la transacción en caso de error
		 tx.Rollback()
		 return err
	}

	// Insertar nuevas filas en la tabla
	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		 // Deshacer la transacción en caso de error
		 tx.Rollback()
		 return err
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		 return err
	}

	fmt.Println("==== Datos CPU Insertados ====")
	return nil
}
