package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"github.com/veritasos/runnel/runnel"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// JSONResponse sent by server
type JSONResponse struct {
	Response string `json:"response"`
}

// Application command from user
type Application struct {
	Cmd string `json:"cmd"`
	Cwd string `json:"cwd"`
}

func command(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	app := new(Application)
	if err := json.NewDecoder(req.Body).Decode(&app); err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Split command into executable and args
	executableAndArgs := strings.Split(app.Cmd, " ")

	client := runnel.NewClient()
	key, err := client.RunCommand(
		executableAndArgs[0], executableAndArgs[1:],
		app.Cwd)

	if err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
	data := &JSONResponse{Response: key}
	if err := json.NewEncoder(res).Encode(data); err != nil {
		panic(err)
	}

}

func stream(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(req)
	timeout, _ := strconv.Atoi(req.FormValue("timeout"))

	client := runnel.NewClient()
	output, err := client.Stream(vars["key"], timeout)

	if err != nil {
		log.Println(err.Error())
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	data := &JSONResponse{Response: output}
	if err := json.NewEncoder(res).Encode(data); err != nil {
		panic(err)
	}
}

func home(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	message := `
	Available endpoints
	http://localhost:9090/command
	http://localhost:9090/stream
	`

	data := &JSONResponse{Response: message}
	if err := json.NewEncoder(res).Encode(data); err != nil {
		panic(err)
	}
}

func main() {
	var ip = flag.String("ip-address", "", "ip:port")
	flag.StringVar(ip, "p", "", "ip:port")
	flag.Parse()

	if *ip == "" {
		log.Fatal("Provide ip address and port, '-p ip:port'")
	}

	log.Println("Started Runnel Server on", *ip)

	router := mux.NewRouter()
	router.HandleFunc("/", home).Methods("GET")
	router.HandleFunc("/command", command).Methods("POST")
	router.HandleFunc("/stream/{key}", stream).Methods("GET")
	http.ListenAndServe(*ip, router)
}
