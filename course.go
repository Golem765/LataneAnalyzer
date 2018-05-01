package main

import (
	"fmt"
	"univ/course/db"
	"github.com/spf13/viper"
	"univ/course/service/user"
	"univ/course/model"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"io/ioutil"
	"encoding/json"
	"univ/course/analyzer"
	"univ/course/service/mentions"
)

var (
	database       db.DB
	repositoryName = db.Users
)

func main() {
	err := readConfig()
	if err != nil {
		fmt.Printf("could not read config: %v", err)
		return
	}

	database, err = connectDatabase()
	if err != nil {
		fmt.Printf("could not establish db connection: %v", err)
		return
	}
	defer database.Close()

	router := mux.NewRouter()
	router.HandleFunc("/go", Analyze).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./frontend/")))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	handler := c.Handler(router)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":8000",
	}

	log.Fatal(srv.ListenAndServe())
}

type Request struct {
	Beta     float64 `json:"beta"`
	Alpha     float64 `json:"alpha"`
	Simulate bool    `json:"simulate"`
	Users    int     `json:"users"`
}

type Response struct {
	Users []model.User `json:"users"`
}

func Analyze(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var req Request

	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}

	var userService userservice.Service
	if req.Simulate {
		userService = userservice.NewFake(req.Users)
	} else {
		userService = userservice.NewService(database, repositoryName)
	}

	beta, alpha := req.Beta, req.Alpha
	users, err := userService.LoadChecked()
	if err != nil {
		fmt.Printf("load users error: %v", err)
	}

	users = mentions.ProcessUserMentions(users, userService)

	users = analyzer.FindPressure(users, beta, alpha)
	users = analyzer.FindInfluence(users)

	resp := Response{
		Users: users,
	}

	json.NewEncoder(w).Encode(resp)
}

func readConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	return viper.ReadInConfig()
}

func connectDatabase() (db.DB, error) {
	dbConfig := db.MgoConfig{
		Url:      viper.GetString("database.url"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
		Name:     viper.GetString("database.name"),
	}

	return db.NewMgo(dbConfig)
}
