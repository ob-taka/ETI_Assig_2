package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Module struct {
	Mod_ID    string `json:"mod_id"`
	Mod_Name  string `json:"mod_name"`
	Mod_Syno  string `json:"mod_syno"`
	Mod_LO    string `json:"mod_lo"`
	Mod_Class string `json:"mod_class"`
	Tutor     string `json:"tutor"`
	Student   string `json:"student"`
	Mod_RnC   string `json:"mod_rnc"`
}

//test
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "welcome")
}

func getMod() []*Module {
	db, err := sql.Open("mysql", "test:password@tcp(db:3306)/Modules")

	//check for err for db connection
	if err != nil {
		panic(err)
	}

	defer db.Close()

	//execute query
	results, err := db.Query("SELECT * FROM module")
	if err != nil {
		panic(err)
	}

	var moudles []*Module
	for results.Next() {
		var m Module
		err = results.Scan(&m.Mod_ID, &m.Mod_Name, &m.Mod_Syno, &m.Mod_LO, &m.Mod_Class, &m.Tutor, &m.Student, &m.Mod_RnC)
		if err != nil {
			panic(err)
		}
		moudles = append(moudles, &m)
	}
	return moudles
}

//user page test
func modpage(w http.ResponseWriter, r *http.Request) {
	mods := getMod()

	json.NewEncoder(w).Encode(mods)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/mod", modpage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
