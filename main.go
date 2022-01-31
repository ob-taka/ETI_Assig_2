package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Module struct {
	Mod_ID    string   `json:"mod_id"`
	Mod_Name  string   `json:"mod_name"`
	Mod_Syno  string   `json:"mod_syno"`
	Mod_LO    string   `json:"mod_lo"`
	Mod_Class []string `json:"mod_class"`
	Tutor     []string `json:"tutor"`
	Student   []string `json:"student"`
	Mod_RnC   string   `json:"mod_rnc"`
}

//test
func connect() *sql.DB {
	db, err := sql.Open("mysql", "test:password@tcp(db:3306)/Modules")

	//check for err for db connection
	if err != nil {
		panic(err)
	}

	return db
}

func getAllMod() []Module {

	db := connect()
	//execute query
	results, err := db.Query("SELECT * FROM module")
	if err != nil {
		panic(err)
	}

	var moudles []Module
	for results.Next() {
		var m Module
		err = results.Scan(&m.Mod_ID, &m.Mod_Name, &m.Mod_Syno, &m.Mod_LO, &m.Mod_RnC)
		if err != nil {
			panic(err)
		}

		m.Tutor = append(m.Tutor, getTutor(m.Mod_ID)...)
		m.Student = append(m.Student, getStudent(m.Mod_ID)...)
		m.Mod_Class = append(m.Mod_Class, getClass(m.Mod_ID)...)
		moudles = append(moudles, m)
	}

	defer db.Close()
	return moudles
}

func getMod(id string) []Module {

	db := connect()
	//execute query
	results, err := db.Query("SELECT * FROM module where moduleCode = ?", id)
	if err != nil {
		panic(err)
	}

	var moudles []Module
	for results.Next() {
		var m Module
		err = results.Scan(&m.Mod_ID, &m.Mod_Name, &m.Mod_Syno, &m.Mod_LO, &m.Mod_RnC)
		if err != nil {
			panic(err)
		}

		m.Tutor = append(m.Tutor, getTutor(m.Mod_ID)...)
		m.Student = append(m.Student, getStudent(m.Mod_ID)...)
		m.Mod_Class = append(m.Mod_Class, getClass(m.Mod_ID)...)
		moudles = append(moudles, m)
	}

	defer db.Close()
	return moudles
}

func getTutor(id string) []string {
	db := connect()
	results, err := db.Query("select tID from tutor where tMoulde = ?", id)
	if err != nil {
		panic(err)
	}

	var tutor []string
	for results.Next() {
		var t string
		err = results.Scan(&t)
		if err != nil {
			panic(err)
		}
		tutor = append(tutor, t)
	}

	defer db.Close()
	return tutor
}

func getStudent(id string) []string {
	db := connect()
	results, err := db.Query("SELECT stuID FROM student WHERE sModuel = ?", id)
	if err != nil {
		panic(err)
	}

	var student []string
	for results.Next() {
		var s string
		err = results.Scan(&s)
		if err != nil {
			panic(err)
		}
		student = append(student, s)
	}
	defer db.Close()
	return student
}

func getClass(id string) []string {
	db := connect()

	results, err := db.Query("SELECT classCode FROM class WHERE mouduleCode = ?", id)
	if err != nil {
		panic(err)
	}

	var class []string
	for results.Next() {
		var c string
		err = results.Scan(&c)
		if err != nil {
			panic(err)
		}
		class = append(class, c)
	}
	defer db.Close()
	return class
}

//user page test
func showAllMod(w http.ResponseWriter, r *http.Request) {

	mod := getAllMod()

	json.NewEncoder(w).Encode(mod)
}

func searchMod(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	mod := getMod(params["id"])

	json.NewEncoder(w).Encode(mod)
}

func main() {

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET", "PUT", "POST"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/api/v1/Module", showAllMod).Methods("GET")
	r.HandleFunc("/api/v1/Module/{id}", searchMod).Methods("GET")
	log.Fatal(http.ListenAndServe(":8171", handlers.CORS(headers, methods, origins)(r)))
}
