# ETI Assignment 2

Docker image link:

> [https://hub.docker.com/repository/docker/panzeo/eti_assig_2](https://hub.docker.com/repository/docker/panzeo/eti_assig_2)
> 
Set up:
```
docker run panzeo/eti_assig_2
```
OR
```
git clone https://github.com/ob-taka/ETI_Assig_2
```
```
docker-compose up -d
```

# API architecture diagram


![Untitled](https://user-images.githubusercontent.com/48742733/152543324-a5391c0b-0953-4ff5-aa5f-d081f59deee8.png)


# Architecture design consideration



## Microservice design

The architecture for this assignment is relatively straight forward, it is broken down to 3 different API services namely, Class, Module, and View all module API. However since I was assign the microservice of viewing all the module that is available, hence, only View all module API and web frontend is done by me, however this service is dependent on the responses of the other mentioned API to fulfill it’s duty.

```jsx
methods := handlers.AllowedMethods([]string{"GET"})
```

Since the view module API only responsible for displaying information, only ***GET*** request are necessary, out of which there are 2 handlers to handle the request from the web frontend. 

- `api/v1/Module`: retrieve all available modules with individual module information including  Module code, Module name, synopsis, learning objective, classes assigned, tutors assigned, enrolled student and ratings and comments.
- `api/v1/Module/{id}`: retrieve specified module with information stated above.

Even though the API depends on the other API to function, it is still fully independent from the other API and server at [localhost](http://localhost) endpoint with a port number of ***8171*** while the frontend Nginx server serve at port 8170 on their onw docker container.



## Nginx

Nginx is a open source web server that store and delivers content for the web frontend created for this project, this allow the browser on the local machine to send http request to the containerized webserver and deliver the content that it retrieve from our API.

**Set up for Nginx:**

In order to properly use the webserver we have to make some changes so that it would not clash with other microservice for this project 

```jsx
COPY ./module.html /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/
```

In the Dockerfile for the frontend server we specified that `module.html` would be copied to the html folder which is where all the html file is located by default. We also have a custom config file that we want to copy and replace the original config file in the directory.

```bash
server {
        # specify the porst to lisent to, default 80
        listen                  8170;
        # specify the default folder that contains the html files
        root                    /usr/share/nginx/html;
        # change the default page from index.html to module.html since it's not the homepage for the overall project
        index                   module.html;
        server_name             localhost;
        client_max_body_size    16m;
        # allow request from external sources in case other services needs it
        add_header 'Access-Control-Allow-Origin' '*';
    }
```

In the code snip we’ve change the default port the service is listening to to port **8170** as specified by the requirement, we also set the index page to point to `module.htm` so that on start up it knows which page to load, so that there won’t be any confusion if other microservice uses index as they html file name. lastly http request form other service to the page is also allowed in case there is a need from the other microservices. 



# Mock data



Since this is a class wide project, there is bound to have a certain level of incompatibility between the microservices, which is why some of the data are being hardcoded as place holder to simulate a success response from the APIs in the case the dependent microservices run into problems.

### Update log: 31 Jan 2022

**Class management:**

> `GET localhost:8041/api/v1/classes/{semester_start_date}`
> 

```bash
startDate :=  "16-01-2022"
url := "localhost:8041/api/v1/classes/" + startDate
```

Sample JSON:

```json
{"SemesterModules":[
		{"ModuleClasses":[
			{"Capacity": 10,
			"ClassCode":"ADB02",
			"Schedule":"17-01-2022 10:00:00",
			"Students": [],
			"Tutor":"T001"
			},
			{"Capacity": 11,
			"ClassCode":"ADB01",
			"Schedule":"17-01-2022 10:00:00",
			"Students":["S004","S005","S006"],
			"Tutor":"T002"
			}],
			"ModuleCode":"CM",
			"ModuleName":"Computing Math"}],
			"SemesterStartDate":"16-01-2022"}
```

the Class API from class management microservice requires a `semester_start_date` variable in the handler to display **Course that are currently active within a given semester (assigned at least one class/student/tutor), t**his is not ideal as when the students first boot up the page ****there isn’t a way to specify a date as required by the API which would result in a page with all the modules with class, student, and tutor field empty, hence it is hardcoded to simulate a successful response from the API

**Module Management:**

```go
url := "http://localhost:8114/query"
jsonBody := map[string]string{"query": "query SearchModules($text: String!){searchModules(text: $text){name, synopsis}}, "variables":{"text": modCode}"}
jsonValue, _ := json.Marshal(jsonBody)
request, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonValue))
```

Sample JSON:

```json
{"data":{
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
	}
```

Based on the designer’s requirement, a JSON body is require when making a HTTP request, this again has to be implemented on a separate microservice  however, as of this update only `name` and `sysnopsis` field is available, hence a sample JSON is used to simulate a complete response from the API



# Backup Database



in the even that the other API and database are not working due to unforeseeable reason, or the sole purpose is to demo the microservice working independently with direct access to the database,  update the code in the `main.go` file found between line 373 to 391

   

```GO
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
```
