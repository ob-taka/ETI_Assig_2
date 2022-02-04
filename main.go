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

//Final module json object to be pick up by frontend
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

//API Json response from Module management microservices
type ModAPI struct {
	Data struct {
		Modules []struct {
			LearningObj string `json:"learningObj"`
			ModuleCode  string `json:"moduleCode"`
			ModuleName  string `json:"moduleName"`
			Synopsis    string `json:"synopsis"`
		} `json:"modules"`
	} `json:"data"`
}

type SearchModAPI struct {
	Data struct {
		SearchModules []struct {
			LearningObj string `json:"learningObj"`
			ModuleCode  string `json:"moduleCode"`
			ModuleName  string `json:"moduleName"`
			Synopsis    string `json:"synopsis"`
		} `json:"searchModules"`
	} `json:"data"`
}

//API Json response form class management microservices
type classAPI struct {
	SemesterModules []struct {
		ModuleClasses []struct {
			Capacity  int64    `json:"Capacity"`
			ClassCode string   `json:"ClassCode"`
			Schedule  string   `json:"Schedule"`
			Students  []string `json:"Students"`
			Tutor     string   `json:"Tutor"`
		} `json:"ModuleClasses"`
		ModuleCode string `json:"ModuleCode"`
		ModuleName string `json:"ModuleName"`
	} `json:"SemesterModules"`
	SemesterStartDate string `json:"SemesterStartDate"`
}

//Custome module class structure to be parse from API response
type modClass struct {
	moduleCode string   `json:"ModuleCode"`
	student    []string `json:"Student"`
	tutor      []string `json:"Tutor"`
	classCode  []string `json:"ClassCode"`
}

//Connect to build in database
//incase API response dont work
func connect() *sql.DB {
	db, err := sql.Open("mysql", "test:password@tcp(db:3306)/Modules")

	//check for err for db connection
	if err != nil {
		panic(err)
	}

	return db
}

//Get all module from the database and map to Module struct
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

//Get all module available from Moudle managements API
func getModAPI() []Module {
	var module []Module
	var temp Module
	var mod ModAPI

	// //as of 31 Jan Module API only provide module name and synopsis
	// url := "http://localhost:8114/query"
	// // required JSON body as per design of the API
	// jsonBody := map[string]string{"query": "query ListModules(){modules(){name, synopsis}}"}
	// jsonValue, _ := json.Marshal(jsonBody)
	// request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	// request.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// response, err := client.Do(request)
	// if err != nil{
	// 	fmt.Println("request failed")
	// 	return module
	// }
	// resBody, _ := ioutil.ReadAll(response.Body)
	// jsonArr := json.Unmarshal(resBody, &mod)

	//assuming the API call works, the body of the API response structure should be as follows:
	data := `{"data":{
		"modules":[
			{"moduleCode": "CM", 
			"moduleName": "Computing Math", 
			"synopsis": "Learn about computing math",
			 "learningObj": "learn CM"
			 },
			 {"moduleCode": "DB",
			  "moduleName": "Data Base", 
			  "synopsis": "Learn about Data Base", 
			  "learningObj": "learn DB"
			  }
			]
		}
	}`
	json.Unmarshal([]byte(data), &mod)
	class := getClassAPI()
	for _, i := range mod.Data.Modules {
		temp.Mod_ID = i.ModuleCode
		temp.Mod_Name = i.ModuleName
		temp.Mod_Syno = i.Synopsis
		temp.Mod_LO = i.LearningObj
		temp.Mod_RnC = "Moudle is fun" // ask 3.9 idiot...
		for _, c := range class {
			if temp.Mod_ID == c.moduleCode {
				temp.Mod_Class = append(temp.Mod_Class, c.classCode...)
				temp.Student = append(temp.Student, c.student...)
				temp.Tutor = append(temp.Tutor, c.tutor...)
			}
		}
		module = append(module, temp)
		temp = Module{}
	}
	return module
}

//Get sepcified from Moudle managements API
func searchModAPI(modCode string) []Module {
	var module []Module
	var temp Module
	var mod SearchModAPI

	// //as of 31 Jan Module API only provide module name and synopsis
	// url := "http://localhost:8114/query"
	// // required JSON body as per design of the API
	// jsonBody := map[string]string{"query": "query SearchModules($text: String!){searchModules(text: $text){name, synopsis}}, "variables":{"text": modCode}"}
	// jsonValue, _ := json.Marshal(jsonBody)
	// request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
	// request.Header.Set("Content-Type", "application/json")
	// client := &http.Client{}
	// response, err := client.Do(reque st)
	// if err != nil{
	// 	fmt.Println("request failed")
	// 	return module
	// }
	// resBody, _ := ioutil.ReadAll(response.Body)
	// jsonArr := json.Unmarshal(resBody, &mod)

	//assuming the API call works, the body of the API response structure should be as follows:
	data := `{"data":{
		"searchModules":[
			{"moduleCode": "CM", 
			"moduleName": "Computing Math", 
			"synopsis": "Learn about computing math",
			 "learningObj": "learn CM"
			 }
			]
		}
	}`
	json.Unmarshal([]byte(data), &mod)
	class := getClassAPI()
	for _, i := range mod.Data.SearchModules {
		temp.Mod_ID = i.ModuleCode
		temp.Mod_Name = i.ModuleName
		temp.Mod_Syno = i.Synopsis
		temp.Mod_RnC = "Moudle is fun"
		for _, c := range class {
			if temp.Mod_ID == c.moduleCode {
				temp.Mod_Class = append(temp.Mod_Class, c.classCode...)
				temp.Student = append(temp.Student, c.student...)
				temp.Tutor = append(temp.Tutor, c.tutor...)
			}
		}
		module = append(module, temp)
		temp = Module{}
	}
	return module
}

//Get specified module from database: Search
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

//Get tutors ID assigned base the module code arguement parse in
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

//Get students ID assigned base on the module code arguement parse in
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

//Get class code assigned base on the module code argument parse in
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

//Get all the classes
func getClassAPI() []modClass {
	/*
	   Class management API it is unable to retrive any infromation without a semester start date
	   hence a date would be hardcoded, idealy the API should response with all current actived
	   modules unless specified
	*/
	var result []modClass
	var class classAPI
	var temp modClass

	// startDate :=  "16-01-2022"
	// url := "localhost:8041/api/v1/classes/" + startDate
	// response, err := http.Get(url)
	// if err != nil {
	// 	fmt.Println("request failed")
	// 	return result
	// }
	// data, _ := ioutil.ReadAll(response.Body)
	// jsonArr := json.Unmarshal([]byte(data), &class)

	//Assuming the API successfully respond the response body should be as follows:
	data := `{"SemesterModules":[
		{"ModuleClasses":[
			{"Capacity":10,
			"ClassCode":"ADB02",
			"Schedule":"17-01-2022 10:00:00",
			"Students": [],
			"Tutor":"T001"
			},
			{"Capacity":11,
			"ClassCode":"ADB01",
			"Schedule":"17-01-2022 10:00:00",
			"Students":["S004","S005","S006"],
			"Tutor":"T002"
			}],
			"ModuleCode":"CM",
			"ModuleName":"Computing Math"}],
			"SemesterStartDate":"16-01-2022"}`
	json.Unmarshal([]byte(data), &class)

	for _, i := range class.SemesterModules {
		temp.moduleCode = i.ModuleCode
		for _, x := range i.ModuleClasses {
			temp.classCode = append(temp.classCode, x.ClassCode)
			temp.student = append(temp.student, x.Students...)
			temp.tutor = append(temp.tutor, x.Tutor)
		}
		result = append(result, temp)
	}

	return result
}

func showAllMod(w http.ResponseWriter, r *http.Request) {
	//Retrive data from Database
	//mod := getAllMod()

	//Retrive data from API response
	mod := getModAPI()

	json.NewEncoder(w).Encode(mod)
}

func searchMod(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	//mod := getMod(params["id"])
	mod := searchModAPI(params["id"])

	json.NewEncoder(w).Encode(mod)
}

func main() {

	r := mux.NewRouter()
	headers := handlers.AllowedHeaders([]string{"X-REQUESTED-With", "Content-Type"})
	methods := handlers.AllowedMethods([]string{"GET"})
	origins := handlers.AllowedOrigins([]string{"*"})
	r.HandleFunc("/api/v1/Module", showAllMod).Methods("GET")
	r.HandleFunc("/api/v1/Module/{id}", searchMod).Methods("GET")
	log.Fatal(http.ListenAndServe(":8171", handlers.CORS(headers, methods, origins)(r)))
}
